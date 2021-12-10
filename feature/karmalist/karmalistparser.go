package karmalist

import (
	"errors"
	"log"

	"github.com/aetimmes/discordgo"
	"github.com/jakevoytko/crbot/model"
)

// Parser parses ?karmalist commands.
type Parser struct{}

// NewParser works as advertised.
func NewParser() *Parser {
	return &Parser{}
}

// GetName returns the named type of this feature.
func (p *Parser) GetName() string {
	return model.CommandNameKarmaList
}

const (
	// MsgHelpList is the help text for ?karmalist
	MsgHelpKarmaList = "Type `?karmalist` to get the URL of a hastebin with all them karma"
)

// HelpText explains how to use ?list.
func (p *Parser) HelpText(command string) (string, error) {
	return MsgHelpKarmaList, nil
}

// Parse parses the given list command.
func (p *Parser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseList called with non-list command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.CommandTypeKarmaList,
	}, nil
}
