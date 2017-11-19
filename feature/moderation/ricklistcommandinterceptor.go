package moderation

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/config"
	"github.com/jakevoytko/crbot/model"
)

// RickListCommandInterceptor asserts that the
type RickListCommandInterceptor struct {
	rickList []model.Snowflake
}

// NewRickListCommandInterceptor returns a new ricklist command interceptor.
func NewRickListCommandInterceptor(config *config.Config) *RickListCommandInterceptor {
	return &RickListCommandInterceptor{
		rickList: config.RickList,
	}
}

// Intercept checks whether the command is forbidden by the ricklist.
func (i *RickListCommandInterceptor) Intercept(command *model.Command, s api.DiscordSession) (*model.Command, error) {
	// Check moderation.
	// RickList
	// - RickListed users can only use ?learn in private channels, without it responding with
	//   a rickroll.
	if channel, err := s.Channel(command.ChannelID.Format()); err == nil {
		isPrivate := channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM
		isAllowed := command.Type == model.Type_Learn || command.Type == model.Type_None
		if isPrivate && !isAllowed {
			for _, ricked := range i.rickList {
				if ricked.Format() == command.Author.ID {
					return &model.Command{
						Type:      model.Type_RickList,
						Author:    nil,
						ChannelID: command.ChannelID,
					}, nil
				}
			}
		}
	}

	return command, nil
}
