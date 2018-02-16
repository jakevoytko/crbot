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

// LearnParser parses ?learn commands.
type LearnParser struct {
	featureRegistry *feature.Registry
	commandMap      stringmap.StringMap
}

// NewLearnParser works as advertised.
func NewLearnParser(featureRegistry *feature.Registry, commandMap stringmap.StringMap) *LearnParser {
	return &LearnParser{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (p *LearnParser) GetName() string {
	return model.Name_Learn
}

// HelpText explains how to use ?learn.
func (p *LearnParser) HelpText(command string) (string, error) {
	return MsgHelpLearn, nil
}

// Parse parses the given learn command.
func (f *LearnParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != f.GetName() {
		log.Fatal("parseLearn called with non-learn command", errors.New("wat"))
	}
	splitContent = util.CollapseWhitespace(splitContent, 1)
	splitContent = util.CollapseWhitespace(splitContent, 2)

	callRegexp := regexp.MustCompile("^[[:alnum:]].*$")
	responseRegexp := regexp.MustCompile("(?s)^[^/?!].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 3 || !callRegexp.MatchString(splitContent[1]) || !responseRegexp.MatchString(splitContent[2]) {
		return &model.Command{
			Type: model.Type_Help,
			Help: &model.HelpData{
				Command: model.Name_Learn,
			},
		}, nil
	}

	// Don't overwrite old or builtin commands.
	has, err := f.commandMap.Has(splitContent[1])
	if err != nil {
		return nil, err
	}
	if has || f.featureRegistry.IsInvokable(splitContent[1]) {
		return &model.Command{
			Type: model.Type_Learn,
			Learn: &model.LearnData{
				CallOpen: false,
				Call:     splitContent[1],
			},
		}, nil
	}

	// Everything is good.
	response := strings.Join(splitContent[2:], " ")
	return &model.Command{
		Type: model.Type_Learn,
		Learn: &model.LearnData{
			CallOpen: true,
			Call:     splitContent[1],
			Response: response,
		},
	}, nil
}
