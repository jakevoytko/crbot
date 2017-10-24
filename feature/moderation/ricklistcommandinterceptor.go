package moderation

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/model"
)

// RickListCommandInterceptor asserts that the
type RickListCommandInterceptor struct {
	rickList []model.Snowflake
}

// NewRickListCommandInterceptor returns a new ricklist command interceptor.
func NewRickListCommandInterceptor(config *app.Config) *RickListCommandInterceptor {
	return &RickListCommandInterceptor{
		rickList: config.RickList,
	}
}

// Intercept checks whether the command is forbidden by the ricklist.
func (i *RickListCommandInterceptor) Intercept(command *model.Command, s api.DiscordSession, m *discordgo.MessageCreate) (*model.Command, error) {
	// Check moderation.
	// RickList
	// - RickListed users can only use ?learn in private channels, without it responding with
	//   a rickroll.
	if channel, err := s.Channel(m.ChannelID); err == nil && (channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM) && command.Type != model.Type_Learn {
		for _, ricked := range i.rickList {
			if ricked.Format() == m.Author.ID {
				return &model.Command{
					Type: model.Type_RickList,
				}, nil
			}
		}
	}

	return command, nil
}
