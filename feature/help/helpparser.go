package help

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
)

// HelpParser parses ?help commands.
type HelpParser struct {
	featureRegistry *feature.Registry
}

// NewHelpParser works as advertised.
func NewHelpParser(featureRegistry *feature.Registry) *HelpParser {
	return &HelpParser{
		featureRegistry: featureRegistry,
	}
}

// GetName returns the named type.
func (p *HelpParser) GetName() string {
	return model.Name_Help
}

const (
	MsgHelpHelp = "Get help for help."
)

// GetHelpText returns the help text.
func (p *HelpParser) HelpText(command string) (string, error) {
	return MsgHelpHelp, nil
}

// Parse parses the given help command.
func (p *HelpParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	splitContent = util.CollapseWhitespace(splitContent, 1)

	userCommand := ""
	if len(splitContent) > 1 {
		userCommand = splitContent[1]
	}

	return &model.Command{
		Type: model.Type_Help,
		Help: &model.HelpData{
			Command: userCommand,
		},
	}, nil
}
