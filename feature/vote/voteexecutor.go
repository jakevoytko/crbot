package vote

import (
	"fmt"
	"strconv"

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
	MsgBroadcastNewVote = "@here %s started a new vote: %s"
)

// Execute uploads the command list to github and pings the gist link in chat.
func (e *VoteExecutor) Execute(s api.DiscordSession, channel string, command *model.Command) {
	ok, err := e.modelHelper.IsVoteActive()
	if err != nil {
		log.Fatal("Error occurred while calling for active vote", err)
	}
	if ok {
		_, err := s.ChannelMessageSend(channel, MsgActiveVote)
		if err != nil {
			log.Fatal("Unable to send vote-already-active message to user", err)
		}
		return
	}

	userID, err := strconv.ParseInt(command.Author.ID, 10 /* base */, 64 /* bitSize */)
	if err != nil {
		log.Info("Error parsing command user ID", err)
		return
	}
	_, err = e.modelHelper.StartNewVote(userID, command.Vote.Message)
	if err != nil {
		log.Fatal("error starting new vote", err)
	}

	broadcastMessage := fmt.Sprintf(MsgBroadcastNewVote, command.Author.Mention(), command.Vote.Message)
	_, err = s.ChannelMessageSend(channel, broadcastMessage)
	if err != nil {
		log.Fatal("Unable to broadcast new message across the channel", err)
	}
}
