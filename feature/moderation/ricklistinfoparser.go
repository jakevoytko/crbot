package moderation

import (
	"errors"
	"log"

	"github.com/jakevoytko/crbot/model"
)

// RickListInfoParser parses ?list commands.
type RickListInfoParser struct{}

// NewRickListInfoParser works as advertised.
func NewRickListInfoParser() *RickListInfoParser {
	return &RickListInfoParser{}
}

// GetName returns the named type of this feature.
func (p *RickListInfoParser) GetName() string {
	return model.Name_RickListInfo
}

const (
	MsgHelpRickListInfo = "Type `?ricklist` to print all users on the moderation rick list. These are whitelisted users who will get rickrolled every time they try to run a private command."
)

// HelpText explains how to use ?ricklist.
func (p *RickListInfoParser) HelpText(command string) (string, error) {
	return MsgHelpRickListInfo, nil
}

// Parse parses the given list command.
func (p *RickListInfoParser) Parse(splitContent []string) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parse ricklist called with non-list command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.Type_RickListInfo,
	}, nil
}
