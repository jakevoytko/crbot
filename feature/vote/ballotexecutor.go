package vote

import (
	"errors"
	"fmt"
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
	MsgVotedAgainst       = "%v voted no"
	MsgVotedInFavor       = "%v voted yes"
	MsgBallotMustBePublic = "Ballots can only be cast in public channels"
)

func (e *BallotExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	userID, err := model.ParseSnowflake(command.Author.ID)
	if err != nil {
		log.Fatal("Error parsing discord user ID", err)
	}

	discordChannel, err := s.Channel(channel.Format())
	if err != nil {
		log.Fatal("This message didn't come from a valid channel", errors.New("wat"))
	}
	if discordChannel.Type == discordgo.ChannelTypeDM || discordChannel.Type == discordgo.ChannelTypeGroupDM {
		s.ChannelMessageSend(channel.Format(), MsgBallotMustBePublic)
		return
	}

	vote, err := e.modelHelper.CastBallot(userID, command.Ballot.InFavor)
	switch err {
	case ErrorNoVoteActive:
		if _, err := s.ChannelMessageSend(channel.Format(), MsgNoActiveVote); err != nil {
			log.Fatal("Unable to send no-active-vote message to user", err)
		}
		return

	case ErrorAlreadyVoted:
		if _, err := s.ChannelMessageSend(channel.Format(), MsgAlreadyVoted); err != nil {
			log.Fatal("Unable to send already voted message to user", err)
		}
		return
	}

	voteMessage := fmt.Sprintf(MsgVotedAgainst, command.Author.Mention())
	if command.Ballot.InFavor {
		voteMessage = fmt.Sprintf(MsgVotedInFavor, command.Author.Mention())
	}

	messages := []string{voteMessage, StatusLine(e.modelHelper.UTCClock, vote)}
	message := strings.Join(messages, "\n")
	if _, err := s.ChannelMessageSend(channel.Format(), message); err != nil {
		log.Info("Failed to send ballot status message", err)
	}
}
