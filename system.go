package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func InitializeRegistry(commandMap StringMap, gist Gist) *FeatureRegistry {
	// Initializing builtin features.
	// TODO(jvoytko): investigate the circularity that emerged to see if there's
	// a better pattern here.
	featureRegistry := NewFeatureRegistry()
	featureRegistry.Register(NewHelpFeature(featureRegistry))
	featureRegistry.Register(NewLearnFeature(featureRegistry, commandMap))
	featureRegistry.Register(NewUnlearnFeature(featureRegistry, commandMap))
	featureRegistry.Register(NewListFeature(featureRegistry, commandMap, gist))
	customFeature := NewCustomFeature(commandMap)
	featureRegistry.Register(customFeature)
	featureRegistry.FallbackFeature = customFeature
	return featureRegistry
}

///////////////////////////////////////////////////////////////////////////////
// Constants
///////////////////////////////////////////////////////////////////////////////

const (
	Type_Custom = iota
	Type_Help
	Type_Learn
	Type_List
	Type_None
	Type_Unlearn
	Type_Unrecognized

	Name_Help    = "?help"
	Name_Learn   = "?learn"
	Name_List    = "?list"
	Name_Unlearn = "?unlearn"
)

///////////////////////////////////////////////////////////////////////////////
// Interfaces
///////////////////////////////////////////////////////////////////////////////

// StringMap stores key/value string pairs. It is always synchronous, but may be
// stored outside the memory space of the program. For instance, in Redis.
type StringMap interface {
	// Has returns whether or not key is present.
	Has(key string) (bool, error)
	// Get returns the given key. Error if key is not present.
	Get(key string) (string, error)
	// Set sets the given key. Allowed to overwrite.
	Set(key, value string) error
	// Delete deletes the given key. Error if key is not present.
	Delete(key string) error
	// GetAll returns every entry as a map.
	GetAll() (map[string]string, error)
}

// DiscordSession is an interface for interacting with Discord within a session
// message handler.
type DiscordSession interface {
	ChannelMessageSend(channel, message string) (*discordgo.Message, error)
}

// Gist is a wrapper around a simple Gist uploader. Returns the URL on success.
type Gist interface {
	Upload(contents string) (string, error)
}

///////////////////////////////////////////////////////////////////////////////
// Controller methods
///////////////////////////////////////////////////////////////////////////////

// getHandleMessage returns the main handler for incoming messages.
func getHandleMessage(commandMap StringMap, featureRegistry *FeatureRegistry) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Never reply to a bot.
		if m.Author.Bot {
			return
		}

		command, err := parseCommand(commandMap, featureRegistry, m.Content)
		if err != nil {
			info("Error parsing command", err)
			return
		}

		feature := featureRegistry.GetFeatureByType(command.Type)
		if feature != nil {
			feature.Execute(s, m.ChannelID, command)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// User message parsing
///////////////////////////////////////////////////////////////////////////////

// HelpData holds data for Help commands.
type HelpData struct {
	Type int
}

type LearnData struct {
	CallOpen bool
	Call     string
	Response string
}

type UnlearnData struct {
	CallOpen bool
	Call     string
}

type CustomData struct {
	Call string
	Args string
}

// TODO(jake): Make this an interface that has only getType(), cast in features.
type Command struct {
	Custom  *CustomData
	Help    *HelpData
	Learn   *LearnData
	Type    int
	Unlearn *UnlearnData
}

// Parses the raw text string from the user. Returns an executable command.
func parseCommand(commandMap StringMap, registry *FeatureRegistry, content string) (*Command, error) {
	if !strings.HasPrefix(content, "?") {
		return &Command{
			Type: Type_None,
		}, nil
	}
	splitContent := strings.Split(content, " ")

	// Parse builtins.
	if feature := registry.GetFeatureByName(splitContent[0]); feature != nil {
		return feature.Parse(splitContent)
	}

	// See if it's a custom command.
	has, err := commandMap.Has(splitContent[0][1:])
	if err != nil {
		info("Error doing custom parsing", err)
		return nil, err
	}
	if has {
		return registry.FallbackFeature.Parse(splitContent)
	}

	// No such command!
	return &Command{
		Type: Type_Unrecognized,
	}, nil
}
