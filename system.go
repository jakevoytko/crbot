package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

func InitializeRegistry(commandMap model.StringMap, gist api.Gist, config *app.Config) *feature.Registry {
	// Initializing builtin features.
	// TODO(jvoytko): investigate the circularity that emerged to see if there's
	// a better pattern here.
	featureRegistry := feature.NewRegistry()
	featureRegistry.Register(NewHelpFeature(featureRegistry))
	featureRegistry.Register(NewLearnFeature(featureRegistry, commandMap))
	featureRegistry.Register(NewListFeature(featureRegistry, commandMap, gist))
	featureRegistry.Register(NewModerationFeature(featureRegistry, config))
	return featureRegistry
}

///////////////////////////////////////////////////////////////////////////////
// Controller methods
///////////////////////////////////////////////////////////////////////////////

// getHandleMessage returns the main handler for incoming messages.
func getHandleMessage(commandMap model.StringMap, featureRegistry *feature.Registry) func(api.DiscordSession, *discordgo.MessageCreate) {
	return func(s api.DiscordSession, m *discordgo.MessageCreate) {
		// Never reply to a bot.
		if m.Author.Bot {
			return
		}

		command, err := parseCommand(commandMap, featureRegistry, m.Content)
		if err != nil {
			log.Info("Error parsing command", err)
			return
		}

		for _, interceptor := range featureRegistry.CommandInterceptors() {
			command, err = interceptor.Intercept(command, s, m)
			if err != nil {
				panic("Ran into error intercepting commands")
			}
		}

		executor := featureRegistry.GetExecutorByType(command.Type)
		if executor != nil {
			executor.Execute(s, m.ChannelID, command)
		}
	}
}

// Parses the raw text string from the user. Returns an executable command.
func parseCommand(commandMap model.StringMap, registry *feature.Registry, content string) (*model.Command, error) {
	if !strings.HasPrefix(content, "?") {
		return &model.Command{
			Type: model.Type_None,
		}, nil
	}
	splitContent := strings.Split(content, " ")

	// Parse builtins.
	if parser := registry.GetParserByName(splitContent[0]); parser != nil {
		return parser.Parse(splitContent)
	}

	// See if it's a custom command.
	has, err := commandMap.Has(splitContent[0][1:])
	if err != nil {
		log.Info("Error doing custom parsing", err)
		return nil, err
	}
	if has {
		return registry.FallbackParser.Parse(splitContent)
	}

	// No such command!
	return &model.Command{
		Type: model.Type_Unrecognized,
	}, nil
}
