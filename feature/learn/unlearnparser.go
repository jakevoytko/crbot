package learn

import (
	"errors"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
)

// UnlearnParser parses ?unlearn commands.
type UnlearnParser struct {
	featureRegistry *feature.Registry
	commandMap      model.StringMap
}

// NewUnlearnParser works as advertised.
func NewUnlearnParser(featureRegistry *feature.Registry, commandMap model.StringMap) *UnlearnParser {
	return &UnlearnParser{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (p *UnlearnParser) GetName() string {
	return model.Name_Unlearn
}

// HelpText returns the help text for ?unlearn.
func (p *UnlearnParser) HelpText(command string) (string, error) {
	return MsgHelpUnlearn, nil
}

// Parse parses the given unlearn command.
func (p *UnlearnParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseUnlearn called with non-unlearn command", errors.New("wat"))
	}

	splitContent = util.CollapseWhitespace(splitContent, 1)

	callRegexp := regexp.MustCompile("^[[:alnum:]].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !callRegexp.MatchString(splitContent[1]) {
		return &model.Command{
			Type: model.Type_Help,
			Help: &model.HelpData{
				Command: model.Name_Unlearn,
			},
		}, nil
	}

	// Only unlearn commands that aren't built-in and exist
	has, err := p.commandMap.Has(splitContent[1])
	if err != nil {
		return nil, err
	}
	if !has || p.featureRegistry.IsInvokable(splitContent[1]) {
		return &model.Command{
			Type: model.Type_Unlearn,
			Unlearn: &model.UnlearnData{
				CallOpen: false,
				Call:     splitContent[1],
			},
		}, nil
	}

	// Everything is good.
	return &model.Command{
		Type: model.Type_Unlearn,
		Unlearn: &model.UnlearnData{
			CallOpen: true,
			Call:     splitContent[1],
		},
	}, nil
}
