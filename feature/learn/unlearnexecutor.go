package learn

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

type UnlearnExecutor struct {
	commandMap model.StringMap
}

func NewUnlearnExecutor(commandMap model.StringMap) *UnlearnExecutor {
	return &UnlearnExecutor{commandMap: commandMap}
}

// GetType returns the type of this feature.
func (f *UnlearnExecutor) GetType() int {
	return model.Type_Unlearn
}

// Execute replies over the given channel indicating successful unlearning, or
// failure to unlearn.
func (e *UnlearnExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	if command.Unlearn == nil {
		log.Fatal("Incorrectly generated unlearn command", errors.New("wat"))
	}

	// Get the current channel and check if we're being asked to unlearn in a
	// private message.
	discordChannel, err := s.Channel(channel.Format())
	if err != nil {
		log.Fatal("This message didn't come from a valid channel", errors.New("wat"))
	}
	if discordChannel.Type == discordgo.ChannelTypeDM || discordChannel.Type == discordgo.ChannelTypeGroupDM {
		s.ChannelMessageSend(channel.Format(), MsgUnlearnMustBePublic)
		return
	}

	if !command.Unlearn.CallOpen {
		s.ChannelMessageSend(channel.Format(), fmt.Sprintf(MsgUnlearnFail, command.Unlearn.Call))
		return
	}

	// Remove the command.
	if has, err := e.commandMap.Has(command.Unlearn.Call); !has || err != nil {
		if has {
			log.Fatal("Tried to unlearn command that doesn't exist: "+command.Unlearn.Call, errors.New("wat"))
		}
		log.Fatal("Error in UnlearnFeature#execute, testing a command", err)
	}
	if err := e.commandMap.Delete(command.Unlearn.Call); err != nil {
		log.Fatal("Unsuccessful unlearning a key; Dying since it might work with a restart", err)
	}

	// Send ack.
	s.ChannelMessageSend(channel.Format(), fmt.Sprintf(MsgUnlearnSuccess, command.Unlearn.Call))
}
