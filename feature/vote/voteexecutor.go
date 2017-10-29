package vote

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

type VoteExecutor struct {
	modelHelper *ModelHelper
}

func NewVoteExecutor(modelHelper *ModelHelper) *VoteExecutor {
	return &VoteExecutor{
		modelHelper: modelHelper,
	}
}

// GetType returns the type of this feature.
func (e *VoteExecutor) GetType() int {
	return model.Type_Vote
}

const (
	MsgActiveVote       = "Cannot start a vote while another is in progress. Type `?votestatus` for more info"
	MsgBroadcastNewVote = "@everyone -- %s started a new vote: %s.\n\nType `?yes` or `?no` to vote. 30 minutes remaining."
	MsgVoteMustBePublic = "Votes can only be started in public channels"
)

// Execute uploads the command list to github and pings the gist link in chat.
func (e *VoteExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	discordChannel, err := s.Channel(channelID.Format())
	if err != nil {
		log.Fatal("This message didn't come from a valid channel", errors.New("wat"))
	}
	if discordChannel.Type == discordgo.ChannelTypeDM || discordChannel.Type == discordgo.ChannelTypeGroupDM {
		s.ChannelMessageSend(channelID.Format(), MsgVoteMustBePublic)
		return
	}

	ok, err := e.modelHelper.IsVoteActive(channelID)
	if err != nil {
		log.Fatal("Error occurred while calling for active vote", err)
	}
	if ok {
		_, err := s.ChannelMessageSend(channelID.Format(), MsgActiveVote)
		if err != nil {
			log.Fatal("Unable to send vote-already-active message to user", err)
		}
		return
	}

	userID, err := model.ParseSnowflake(command.Author.ID)
	if err != nil {
		log.Info("Error parsing command user ID", err)
		return
	}
	_, err = e.modelHelper.StartNewVote(channelID, userID, command.Vote.Message)
	if err != nil {
		log.Fatal("error starting new vote", err)
	}

	broadcastMessage := fmt.Sprintf(MsgBroadcastNewVote, command.Author.Mention(), command.Vote.Message)
	_, err = s.ChannelMessageSend(channelID.Format(), broadcastMessage)
	if err != nil {
		log.Fatal("Unable to broadcast new message across the channel", err)
	}
}
