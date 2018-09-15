package learn

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// Feature allows crbot to learn new calls and responses
type Feature struct {
	featureRegistry *feature.Registry
	commandMap      stringmap.StringMap
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry, commandMap stringmap.StringMap) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// Parsers gets the learn feature parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{
		NewCustomLearnParser(f.featureRegistry, f.commandMap),
		NewUnlearnParser(f.featureRegistry, f.commandMap),
	}
}

// CommandInterceptors returns nothing.
func (f *Feature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns the custom parser, to recognize custom ? commands. It
// should be the only fallback parser in the project.
func (f *Feature) FallbackParser() feature.Parser {
	return NewCustomParser(f.commandMap)
}

// Executors returns the executors for the ?learn feature.
func (f *Feature) Executors() []feature.Executor {
	return []feature.Executor{
		NewCustomLearnExecutor(f.commandMap),
		NewUnlearnExecutor(f.commandMap),
		NewCustomExecutor(f.commandMap),
	}
}

// OnInitialLoad does nothing.
func (f *Feature) OnInitialLoad(s api.DiscordSession) error { return nil }

///////////////////////////////////////////////////////////////////////////////
// Messages
///////////////////////////////////////////////////////////////////////////////

const (
	// MsgCustomNeedsArgs is a user-visible string asking the user for command args.
	MsgCustomNeedsArgs = "This command takes args. Please type `?command <more text>` instead of `?command`"
	// MsgHelpLearn is the help text for ?learn
	MsgHelpLearn = "Type `?learn <call> <the response the bot should read>`. When you type `?call`, the bot will reply with the response.\n\nThe first character of the call must be alphanumeric, and the first character of the response must not begin with /, ?, or !\n\nUse $1 in the response to substitute all arguments"
	// MsgHelpUnlearn is the help text for ?unlearn
	MsgHelpUnlearn = "Type `?unlearn <call>` to forget a user-defined command."
	// MsgLearnFail indicates that the user tried to overwrite a learned command
	MsgLearnFail = "I already know ?%s"
	// MsgLearnSuccess indicates that the bot learned the command
	MsgLearnSuccess = "Learned about %s"
	// MsgUnlearnFail indicates that the user attempted to unlearn an unlearnable command
	MsgUnlearnFail = "I can't unlearn `?%s`"
	// MsgUnlearnMustBePublic indicates that the user tried to unlearn in a private channel
	MsgUnlearnMustBePublic = "I can't unlearn in a private message."
	// MsgUnlearnSuccess indicates the bot deleted the given learn
	MsgUnlearnSuccess = "Forgot about %s"
)
