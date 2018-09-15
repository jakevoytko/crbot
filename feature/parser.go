package feature

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/model"
)

// Parser is used to multiplex on builtin ?* commands, and ensure that the
// commands are correctly formatted.
type Parser interface {
	// The user-facing name of the command. Must be unique.
	GetName() string
	// Parses the given split command line.
	Parse([]string, *discordgo.MessageCreate) (*model.Command, error)
	// The user-facing help text for the given name. `command` is passed as an arg
	// so fallback parsers can provide custom help text.
	HelpText(command string) (string, error)
}
