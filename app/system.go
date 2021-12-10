package app

import (
	"fmt"
	"strings"

	"github.com/aetimmes/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/config"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/feature/factsphere"
	"github.com/jakevoytko/crbot/feature/help"
	"github.com/jakevoytko/crbot/feature/karma"
	"github.com/jakevoytko/crbot/feature/karmalist"
	"github.com/jakevoytko/crbot/feature/learn"
	"github.com/jakevoytko/crbot/feature/list"
	"github.com/jakevoytko/crbot/feature/moderation"
	"github.com/jakevoytko/crbot/feature/vote"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// InitializeRegistry registers all of the features of the chatbot with the registry and returns the initialized
// registry.
func InitializeRegistry(
	commandMap stringmap.StringMap,
	karmaMap stringmap.StringMap,
	voteMap stringmap.StringMap,
	gist api.Gist,
	config *config.Config,
	clock model.UTCClock,
	timer model.UTCTimer,
	commandChannel chan<- *model.Command) *feature.Registry {

	// Initializing builtin features.
	// TODO(jvoytko): investigate the circularity that emerged to see if there's
	// a better pattern here.
	featureRegistry := feature.NewRegistry()
	allFeatures := []feature.Feature{
		factsphere.NewFeature(featureRegistry),
		help.NewFeature(featureRegistry),
		karma.NewFeature(featureRegistry, karmaMap),
		karmalist.NewFeature(featureRegistry, karmaMap, gist),
		learn.NewFeature(featureRegistry, commandMap),
		list.NewFeature(featureRegistry, commandMap, gist),
		moderation.NewFeature(featureRegistry, config),
		vote.NewFeature(featureRegistry, voteMap, clock, timer, commandChannel),
	}

	for _, f := range allFeatures {
		err := featureRegistry.Register(f)
		if err != nil {
			panic(err)
		}
	}
	return featureRegistry
}

///////////////////////////////////////////////////////////////////////////////
// Controller methods
///////////////////////////////////////////////////////////////////////////////

const (
	// MsgPublicOnly is a user-visible string for commands that cannot be executed in a private channels.
	MsgPublicOnly = "Cannot execute `%s` in a private channel"
)

// HandleCommands pops commands off the command channel and attempts to dispatch them to a command executor.
func HandleCommands(featureRegistry *feature.Registry, s api.DiscordSession, commandChannel <-chan *model.Command) {
	for command := range commandChannel {
		var err error // so I don't have to use := in the intercept() call
		for _, interceptor := range featureRegistry.CommandInterceptors() {
			command, err = interceptor.Intercept(command, s)
			if err != nil {
				panic("Ran into error intercepting commands")
			}
		}

		executor := featureRegistry.GetExecutorByType(command.Type)
		if executor != nil {
			if err != nil {
				log.Fatal("Error parsing snowflake", err)
			}
			discordChannel, err := s.Channel(command.ChannelID.Format())
			if err != nil {
				log.Info("Error retrieving channel from discord for command executor", err)
				continue
			}
			if (discordChannel.Type == discordgo.ChannelTypeDM || discordChannel.Type == discordgo.ChannelTypeGroupDM) && executor.PublicOnly() {
				s.ChannelMessageSend(command.ChannelID.Format(), fmt.Sprintf(MsgPublicOnly, command.OriginalName))
				continue
			}
			executor.Execute(s, command.ChannelID, command)
		}
	}
}

// GetHandleMessage returns the main handler for incoming messages.
func GetHandleMessage(commandMap stringmap.StringMap, featureRegistry *feature.Registry, commandChannel chan<- *model.Command) func(api.DiscordSession, *discordgo.MessageCreate) {
	return func(s api.DiscordSession, m *discordgo.MessageCreate) {
		// Never reply to a bot.
		if m.Author.Bot {
			return
		}

		command, err := parseCommand(commandMap, featureRegistry, m)
		if err != nil {
			log.Info("Error parsing command", err)
			return
		}
		command.Author = m.Author
		channelID, err := model.ParseSnowflake(m.ChannelID)
		if err != nil {
			log.Info("Error parsing channel ID", err)
			return
		}
		command.ChannelID = channelID

		commandChannel <- command
	}
}

// Parses the raw text string from the user. Returns an executable command.
func parseCommand(commandMap stringmap.StringMap, registry *feature.Registry, m *discordgo.MessageCreate) (*model.Command, error) {
	content := m.Content
	if !strings.HasPrefix(content, "?") {
		return &model.Command{
			Type: model.CommandTypeNone,
		}, nil
	}
	splitContent := strings.Split(content, " ")

	// Parse builtins.
	if parser := registry.GetParserByName(splitContent[0]); parser != nil {
		command, err := parser.Parse(splitContent, m)
		command.OriginalName = splitContent[0]
		return command, err
	}

	// See if it's a custom command.
	has, err := commandMap.Has(splitContent[0][1:])
	if err != nil {
		log.Info("Error doing custom parsing", err)
		return nil, err
	}
	if has {
		return registry.FallbackParser.Parse(splitContent, m)
	}

	// No such command!
	return &model.Command{
		Type: model.CommandTypeUnrecognized,
	}, nil
}
