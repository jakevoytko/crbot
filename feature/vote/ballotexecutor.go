package vote

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

type BallotExecutor struct {
	modelHelper *ModelHelper
}

func NewBallotExecutor(modelHelper *ModelHelper) *BallotExecutor {
	return &BallotExecutor{
		modelHelper: modelHelper,
	}
}

// GetType returns the type of this feature.
func (e *BallotExecutor) GetType() int {
	return model.Type_VoteBallot
}

const (
	MsgAlreadyVoted       = "You already voted"
	MsgVotedAgainst       = "You voted no"
	MsgVotedInFavor       = "You voted yes"
	MsgBallotMustBePublic = "Ballots can only be cast in public channels"
)

func (e *BallotExecutor) Execute(s api.DiscordSession, channel string, command *model.Command) {
	userID, err := strconv.ParseInt(command.Author.ID, 10, 64)
	if err != nil {
		log.Fatal("Error parsing discord user ID", err)
	}

	discordChannel, err := s.Channel(channel)
	if err != nil {
		log.Fatal("This message didn't come from a valid channel", errors.New("wat"))
	}
	if discordChannel.Type == discordgo.ChannelTypeDM || discordChannel.Type == discordgo.ChannelTypeGroupDM {
		s.ChannelMessageSend(channel, MsgBallotMustBePublic)
		return
	}

	vote, err := e.modelHelper.CastBallot(userID, command.Ballot.InFavor)
	switch err {
	case ErrorNoVoteActive:
		if _, err := s.ChannelMessageSend(channel, MsgNoActiveVote); err != nil {
			log.Fatal("Unable to send no-active-vote message to user", err)
		}
		return

	case ErrorAlreadyVoted:
		if _, err := s.ChannelMessageSend(channel, MsgAlreadyVoted); err != nil {
			log.Fatal("Unable to send already voted message to user", err)
		}
		return
	}

	voteMessage := MsgVotedAgainst
	if command.Ballot.InFavor {
		voteMessage = MsgVotedInFavor
	}

	messages := []string{voteMessage, StatusLine(vote)}
	message := strings.Join(messages, "\n")
	if _, err := s.ChannelMessageSend(channel, message); err != nil {
		log.Info("Failed to send ballot status message", err)
	}
}
