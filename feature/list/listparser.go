package list

import (
	"errors"
	"log"

	"github.com/jakevoytko/crbot/model"
)

// ListParser parses ?list commands.
type ListParser struct{}

// NewListParser works as advertised.
func NewListParser() *ListParser {
	return &ListParser{}
}

// GetName returns the named type of this feature.
func (p *ListParser) GetName() string {
	return model.Name_List
}

const (
	MsgHelpList = "Type `?list` to get the URL of a Gist with all builtin and learned commands"
)

// HelpText explains how to use ?list.
func (p *ListParser) HelpText(command string) (string, error) {
	return MsgHelpList, nil
}

// Parse parses the given list command.
func (p *ListParser) Parse(splitContent []string) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseList called with non-list command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.Type_List,
	}, nil
}
