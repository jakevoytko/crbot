package karma

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/util"
)

// KarmaParser parses ?++ and ?-- commands
type KarmaParser struct {
	// The message that the parser looks for.
	Message string
	// Whether this message is incrementing or decrementing karma
	Increment bool
}

// NewKarmaParser works as advertised.
func NewKarmaParser(message string, increment bool) *KarmaParser {
	return &KarmaParser{
		Message:   message,
		Increment: increment,
	}
}

// GetName returns the named type.
func (p *KarmaParser) GetName() string {
	return p.Message
}

const (
	MsgHelpKarmaIncrement = "Type ?++ <target> to add a single unit of karma to a target's karma score"
	MsgHelpKarmaDecrement = "Type ?-- <target> to deduct a single unit of karma from a target's karma score"
)

// HelpText returns the help text.
func (p *KarmaParser) HelpText(command string) (string, error) {
	if p.Increment {
		return MsgHelpKarmaIncrement, nil
	}
	return MsgHelpKarmaDecrement, nil
}

// Parse parses the given karma command.
func (p *KarmaParser) Parse(splitContent []string) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("KarmaParser.Parse called with non-list command", errors.New("wat"))
	}

	splitContent = util.CollapseWhitespace(splitContent, 1)
	targetRegexp := regexp.MustCompile("^[[:alnum:]].*$")

	// Resolve the target of the karma, removing leading @ and anything after #
	var target string
	if len(splitContent) > 1 {
		target = strings.Split(strings.TrimPrefix(splitContent[1], "@"), "#")[0]
	}

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !targetRegexp.MatchString(target) {
		if p.Increment {
			return &model.Command{
				Type: model.Type_Help,
				Help: &model.HelpData{
					Command: model.Name_KarmaIncrement,
				},
			}, nil
		}
		return &model.Command{
			Type: model.Type_Help,
			Help: &model.HelpData{
				Command: model.Name_KarmaDecrement,
			},
		}, nil
	}

	return &model.Command{
		Type: model.Type_Karma,
		Karma: &model.KarmaData{
			Increment: p.Increment,
			Target:    target,
		},
	}, nil
}
