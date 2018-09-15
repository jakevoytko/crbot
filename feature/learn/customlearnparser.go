package learn

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// CustomLearnParser parses ?learn commands.
type CustomLearnParser struct {
	featureRegistry *feature.Registry
	commandMap      stringmap.StringMap
}

// NewCustomLearnParser works as advertised.
func NewCustomLearnParser(featureRegistry *feature.Registry, commandMap stringmap.StringMap) *CustomLearnParser {
	return &CustomLearnParser{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (p *CustomLearnParser) GetName() string {
	return model.CommandNameLearn
}

// HelpText explains how to use ?learn.
func (p *CustomLearnParser) HelpText(command string) (string, error) {
	return MsgHelpLearn, nil
}

// Parse parses the given learn command.
func (p *CustomLearnParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseLearn called with non-learn command", errors.New("wat"))
	}
	splitContent = util.CollapseWhitespace(splitContent, 1)
	splitContent = util.CollapseWhitespace(splitContent, 2)

	callRegexp := regexp.MustCompile("^[[:alnum:]].*$")
	responseRegexp := regexp.MustCompile("(?s)^[^/?!].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 3 || !callRegexp.MatchString(splitContent[1]) || !responseRegexp.MatchString(splitContent[2]) {
		return &model.Command{
			Type: model.CommandTypeHelp,
			Help: &model.HelpData{
				Command: model.CommandNameLearn,
			},
		}, nil
	}

	// Don't overwrite old or builtin commands.
	has, err := p.commandMap.Has(splitContent[1])
	if err != nil {
		return nil, err
	}
	if has || p.featureRegistry.IsInvokable(splitContent[1]) {
		return &model.Command{
			Type: model.CommandTypeLearn,
			Learn: &model.LearnData{
				CallOpen: false,
				Call:     splitContent[1],
			},
		}, nil
	}

	// Everything is good.
	response := strings.Join(splitContent[2:], " ")
	return &model.Command{
		Type: model.CommandTypeLearn,
		Learn: &model.LearnData{
			CallOpen: true,
			Call:     splitContent[1],
			Response: response,
		},
	}, nil
}
