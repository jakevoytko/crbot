package vote

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/model"
)

// StatusParser parses ?votestatus commands
type StatusParser struct {
}

// NewStatusParser works as advertised.
func NewStatusParser() *StatusParser {
	return &StatusParser{}
}

// GetName returns the named type.
func (p *StatusParser) GetName() string {
	return model.CommandNameVoteStatus
}

const (
	// MsgHelpStatus is the help text for ?votestatus
	MsgHelpStatus = "Prints the status of the current vote, or a message indicating that no vote is active"
)

// HelpText returns the help text.
func (p *StatusParser) HelpText(command string) (string, error) {
	return MsgHelpStatus, nil
}

// Parse parses the given list command.
func (p *StatusParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseVoteStatus called with non-list command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.CommandTypeVoteStatus,
	}, nil
}
