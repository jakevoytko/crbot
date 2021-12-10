package help

import (
	"errors"

	"github.com/aetimmes/discordgo"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
)

// Parser parses ?help commands.
type Parser struct {
	featureRegistry *feature.Registry
}

// NewParser works as advertised.
func NewParser(featureRegistry *feature.Registry) *Parser {
	return &Parser{
		featureRegistry: featureRegistry,
	}
}

// GetName returns the named type.
func (p *Parser) GetName() string {
	return model.CommandNameHelp
}

const (
	// MsgHelpHelp prints the results of ?help help
	MsgHelpHelp = "Get help for help."
)

// HelpText returns the help text.
func (p *Parser) HelpText(command string) (string, error) {
	return MsgHelpHelp, nil
}

// Parse parses the given help command.
func (p *Parser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	splitContent = util.CollapseWhitespace(splitContent, 1)

	userCommand := ""
	if len(splitContent) > 1 {
		userCommand = splitContent[1]
	}

	return &model.Command{
		Type: model.CommandTypeHelp,
		Help: &model.HelpData{
			Command: userCommand,
		},
	}, nil
}
