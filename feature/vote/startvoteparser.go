package vote

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
)

const (
	// MsgHelpVote is the help text for ?vote
	MsgHelpVote = "Type `?vote <message>` to call a yes/no vote on the given message. 5 votes must be cast before the vote can pass. The first character of the message must be alphanumeric.\n\nExample: `?vote are pirates better than ninjas?`"
)

// StartVoteParser parses ?vote commands.
type StartVoteParser struct {
}

// NewStartVoteParser works as advertised.
func NewStartVoteParser() *StartVoteParser {
	return &StartVoteParser{}
}

// GetName returns the named type of this feature.
func (p *StartVoteParser) GetName() string {
	return model.CommandNameVote
}

// HelpText explains how to use ?vote.
func (p *StartVoteParser) HelpText(command string) (string, error) {
	return MsgHelpVote, nil
}

// Parse parses the given vote command.
func (p *StartVoteParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseVote called with non vote command", errors.New("wat"))
	}
	// The command is everything at/after the first word.
	splitContent = util.CollapseWhitespace(splitContent, 1)

	voteRegexp := regexp.MustCompile("^[[:alnum:]].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !voteRegexp.MatchString(splitContent[1]) {
		return &model.Command{
			Type: model.CommandTypeHelp,
			Help: &model.HelpData{
				Command: model.CommandNameVote,
			},
		}, nil
	}

	message := strings.Join(splitContent[1:], " ")
	return &model.Command{
		Type: model.CommandTypeVote,
		Vote: &model.VoteData{
			Message: message,
		},
	}, nil
}
