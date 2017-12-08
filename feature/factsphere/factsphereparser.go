package factsphere

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/model"
)

// FactSphereParser parses ?factsphere commands.
type FactSphereParser struct{}

// NewFactSphereParser works as advertised.
func NewFactSphereParser() *FactSphereParser {
	return &FactSphereParser{}
}

// GetName returns the named type of this feature.
func (p *FactSphereParser) GetName() string {
	return model.Name_FactSphere
}

const (
	MsgHelpFactSphere = "Type `?factsphere` to get a random fact, which may or may not be true."
)

// HelpText explains how to use ?list.
func (p *FactSphereParser) HelpText(command string) (string, error) {
	return MsgHelpFactSphere, nil
}

// Parse parses the given factsphere command.
func (p *FactSphereParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("FactSphereParser.Parse called with non-factsphere command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.Type_FactSphere,
	}, nil
}
