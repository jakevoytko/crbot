package factsphere

import (
	"errors"
	"log"

	"github.com/aetimmes/discordgo"
	"github.com/jakevoytko/crbot/model"
)

// Parser parses ?factsphere commands.
type Parser struct{}

// NewFactSphereParser works as advertised.
func NewFactSphereParser() *Parser {
	return &Parser{}
}

// GetName returns the named type of this feature.
func (p *Parser) GetName() string {
	return model.CommandNameFactSphere
}

const (
	// MsgHelpFactSphere is the help text for ?factsphere.
	MsgHelpFactSphere = "Type `?factsphere` to get a random fact, which may or may not be true."
)

// HelpText explains how to use ?list.
func (p *Parser) HelpText(command string) (string, error) {
	return MsgHelpFactSphere, nil
}

// Parse parses the given factsphere command.
func (p *Parser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("FactSphereParser.Parse called with non-factsphere command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.CommandTypeFactSphere,
	}, nil
}
