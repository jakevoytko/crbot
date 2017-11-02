package feature

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/model"
)

// CommandInterceptor allows features to examine and replace the parsed
// command. For instance, moderation and command ACLs can be implemented this
// way.
type CommandInterceptor interface {
	Intercept(*model.Command, api.DiscordSession) (*model.Command, error)
}
