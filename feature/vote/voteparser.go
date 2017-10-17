package vote

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
)

const (
	MsgHelpVote = "Type `?vote <message>` to call a yes/no vote on the given message. Needs 5 yes votes to pass. The first character of the message must be alphanumeric.\n\nExample: `?vote are pirates better than ninjas?`"
)

// VoteParser parses ?vote commands.
type VoteParser struct {
}

// NewVoteParser works as advertised.
func NewVoteParser() *VoteParser {
	return &VoteParser{}
}

// GetName returns the named type of this feature.
func (p *VoteParser) GetName() string {
	return model.Name_Vote
}

// HelpText explains how to use ?vote.
func (p *VoteParser) HelpText(command string) (string, error) {
	return MsgHelpVote, nil
}

// Parse parses the given vote command.
func (f *VoteParser) Parse(splitContent []string) (*model.Command, error) {
	if splitContent[0] != f.GetName() {
		log.Fatal("parseVote called with non vote command", errors.New("wat"))
	}
	// The command is everything at/after the first word.
	splitContent = util.CollapseWhitespace(splitContent, 1)

	voteRegexp := regexp.MustCompile("^[[:alnum:]].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !voteRegexp.MatchString(splitContent[1]) {
		return &model.Command{
			Type: model.Type_Help,
			Help: &model.HelpData{
				Command: model.Name_Vote,
			},
		}, nil
	}

	message := strings.Join(splitContent[1:], " ")
	return &model.Command{
		Type: model.Type_Vote,
		Vote: &model.VoteData{
			Message: message,
		},
	}, nil
}
