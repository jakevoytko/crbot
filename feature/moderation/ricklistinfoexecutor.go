package moderation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// RickListInfoExecutor prints a rick roll.
type RickListInfoExecutor struct {
	config *app.Config
}

// NewRickListInfoExecutor works as advertised.
func NewRickListInfoExecutor(config *app.Config) *RickListInfoExecutor {
	return &RickListInfoExecutor{
		config: config,
	}
}

// GetType returns the type.
func (e *RickListInfoExecutor) GetType() int {
	return model.Type_RickListInfo
}

const (
	MsgRickListEmpty = "Nobody is on the ricklist."
	MsgRickListUsers = "On the Rick list: "
)

// Execute replies over the given channel with a rick roll.
func (e *RickListInfoExecutor) Execute(s api.DiscordSession, channel string, command *model.Command) {
	if len(e.config.RickList) == 0 {
		if _, err := s.ChannelMessageSend(channel, MsgRickListEmpty); err != nil {
			log.Info("Failed to send ricklist message", err)
		}
		return
	}

	users := make([]string, 0, len(e.config.RickList))
	for _, ricklisted := range e.config.RickList {
		ricklistedFormat := strconv.FormatInt(ricklisted, 10)
		user, err := s.User(ricklistedFormat)
		if err != nil {
			log.Info(fmt.Sprintf("Unable to get info for user %s", ricklisted), err)
			users = append(users, ricklistedFormat)
			continue
		}
		users = append(users, "@"+user.Username)
	}

	finalString := MsgRickListUsers + "[" + strings.Join(users, ", ") + "]"

	if _, err := s.ChannelMessageSend(channel, finalString); err != nil {
		log.Info("Failed to send ricklist message", err)
	}
}
