package feature

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/model"
)

// Executor actually performs the given action based on the given command.
type Executor interface {
	// The command type to execute.
	GetType() int
	// Execute the given command, for the session and channel name provided.
	Execute(api.DiscordSession, model.Snowflake, *model.Command)
	// Whether the command cannot be executed in private channels.
	PublicOnly() bool
}
