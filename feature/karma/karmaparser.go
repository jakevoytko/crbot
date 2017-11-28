package karma

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
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

// The target is mentioned in plaintext.
var directMentionRegexp *regexp.Regexp = regexp.MustCompile("^[[:alnum:]].*$")

// The target is an embedded entity that needs to be looked up in the message.
var entityRegexp *regexp.Regexp = regexp.MustCompile("^<@!?([[:digit:]]+)>$")

// Parse parses the given karma command.
func (p *KarmaParser) Parse(splitContent []string, m *discordgo.MessageCreate) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("KarmaParser.Parse called with non-list command", errors.New("wat"))
	}

	splitContent = util.CollapseWhitespace(splitContent, 1)

	var target string

	if len(splitContent) > 1 {
		// First, test to see if there is an embedded entity that can be looked up.
		entityMatch := entityRegexp.FindStringSubmatch(splitContent[1])
		if len(entityMatch) == 2 {
			// Look up the ID in mentions in the original message.
			id := entityMatch[1]
			for _, mention := range m.Mentions {
				if mention.ID == id {
					target = mention.Username
				}
			}
		}

		// If not, try to trim and match what's left.
		if len(target) == 0 {
			trimmedTarget := strings.Split(strings.TrimPrefix(splitContent[1], "@"), "#")[0]
			if directMentionRegexp.MatchString(trimmedTarget) {
				target = trimmedTarget
			}
		}
	}

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || len(target) == 0 {
		commandType := model.Name_KarmaDecrement
		if p.Increment {
			commandType = model.Name_KarmaIncrement
		}
		return &model.Command{
			Type: model.Type_Help,
			Help: &model.HelpData{
				Command: commandType,
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
