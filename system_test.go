package main

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/util"
)

// TestRunner is a helper that executes messages incrementally, and asserts that
// the global state is always what is expected.
type TestRunner struct {
	// Test object
	T *testing.T

	// State
	GistsCount           int
	DiscordMessagesCount int
	Learns               map[string]*Learn

	// Fakes
	CustomMap      *util.InMemoryStringMap
	Gist           *util.InMemoryGist
	DiscordSession *util.InMemoryDiscordSession

	// Real objects
	FeatureRegistry *FeatureRegistry

	// Controllers under test
	Handler func(DiscordSession, *discordgo.MessageCreate)
}

func (r *TestRunner) AssertState() {
	// Assert counts.
	assertNumCommands(r.T, r.CustomMap, len(r.Learns))
	assertNumGists(r.T, r.Gist, r.GistsCount)
	assertNumDiscordMessages(r.T, r.DiscordSession, r.DiscordMessagesCount)

	// Assert command map state.
	for _, learn := range r.Learns {
		assertCommand(r.T, r.CustomMap, learn.Call, learn.Response)
	}
}

func (r *TestRunner) SendMessage(channel, message, expectedResponse string) {
	sendMessage(r.T, r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	assertNewMessages(r.T, r.DiscordSession,
		[]*util.Message{util.NewMessage(channel, expectedResponse)})
	r.AssertState()
}

func (r *TestRunner) SendLearnMessage(channel, message string, learn *Learn) {
	sendMessage(r.T, r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	r.Learns[learn.Call] = learn
	assertNewMessages(r.T, r.DiscordSession,
		[]*util.Message{util.NewMessage(channel, fmt.Sprintf(MsgLearnSuccess, learn.Call))})
	r.AssertState()
}

func (r *TestRunner) SendUnlearnMessage(channel, message string, call string) {
	sendMessage(r.T, r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	delete(r.Learns, call)
	assertNewMessages(r.T, r.DiscordSession,
		[]*util.Message{util.NewMessage(channel, fmt.Sprintf(MsgUnlearnSuccess, call))})
	r.AssertState()
}

func (r *TestRunner) SendMessageWithoutResponse(channel, message string) {
	sendMessage(r.T, r.DiscordSession, r.Handler, channel, message)
	assertNewMessages(r.T, r.DiscordSession, []*util.Message{})
	r.AssertState()
}

func (r *TestRunner) SendListMessage(channel string) {
	sendMessage(r.T, r.DiscordSession, r.Handler, channel, "?list")
	r.DiscordMessagesCount++
	r.GistsCount++
	assertNewMessages(r.T, r.DiscordSession, []*util.Message{util.NewMessage(channel, "The list of commands is here: https://www.example.com/success")})

	// Assert gist state. Cannot be in AssertState because this would fail at the
	// next learn or unlearn.
	// TODO(jake): Remove duplication between this and listfeature. Maybe just assert number of lines?
	if r.GistsCount > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("List of builtins:\n - ?help: ")
		buffer.WriteString(MsgHelpHelp)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?learn: ")
		buffer.WriteString(MsgHelpLearn)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?list: ")
		buffer.WriteString(MsgHelpList)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?unlearn: ")
		buffer.WriteString(MsgHelpUnlearn)
		buffer.WriteString("\n\n")

		buffer.WriteString("List of learned commands:\n")

		all, _ := r.CustomMap.GetAll()
		custom := make([]string, 0, len(all))
		for name := range all {
			custom = append(custom, name)
		}
		sort.Strings(custom)
		for _, name := range custom {
			buffer.WriteString(" - ?")
			buffer.WriteString(name)
			if strings.Contains(all[name], "$1") {
				buffer.WriteString(" <args>")
			}
			buffer.WriteString("\n")
		}

		generated := buffer.String()
		actual := r.Gist.Messages[len(r.Gist.Messages)-1]
		if generated != actual {
			r.T.Fatalf(fmt.Sprintf("Gist failure, got `%v` expected `%v`", actual, generated))
		}
	}

	r.AssertState()
}

func Test_Integration(t *testing.T) {
	// Initialize fakes.
	customMap := util.NewInMemoryStringMap()
	gist := util.NewInMemoryGist()
	discordSession := util.NewInMemoryDiscordSession()
	discordSession.SetChannel(&discordgo.Channel{
		ID:        "channel",
		IsPrivate: false,
	})
	discordSession.SetChannel(&discordgo.Channel{
		ID:        "literally anything else",
		IsPrivate: true,
	})

	registry := InitializeRegistry(customMap, gist)
	runner := &TestRunner{
		T:                    t,
		Learns:               map[string]*Learn{},
		GistsCount:           0,
		DiscordMessagesCount: 0,
		CustomMap:            customMap,
		Gist:                 gist,
		DiscordSession:       discordSession,
		FeatureRegistry:      registry,
		Handler:              getHandleMessage(customMap, registry),
	}

	// Assert initial state.
	runner.AssertState()

	// Commands that should never return a response.
	runner.SendMessageWithoutResponse("channel", "?")
	runner.SendMessageWithoutResponse("channel", "!")
	runner.SendMessageWithoutResponse("channel", ".")
	runner.SendMessageWithoutResponse("channel", "")
	runner.SendMessageWithoutResponse("channel", "!help")
	runner.SendMessageWithoutResponse("channel", "help")
	runner.SendMessageWithoutResponse("channel", ".help")

	// Test ?list. ?list tests will be interspersed through the learn examples
	// below, since learn and unlearn interact with it.
	runner.SendListMessage("channel")

	// Basic learn responses.
	// Wrong call format
	runner.SendMessage("channel", "?learn", MsgHelpLearn)
	runner.SendMessage("channel", "?learn test", MsgHelpLearn)
	runner.SendMessage("channel", "?learn ?call response", MsgHelpLearn)
	runner.SendMessage("channel", "?learn !call response", MsgHelpLearn)
	runner.SendMessage("channel", "?learn /call response", MsgHelpLearn)
	runner.SendMessage("channel", "?learn ", MsgHelpLearn)
	runner.SendMessage("channel", "?learn multi\nline\ncall response", MsgHelpLearn)
	// Wrong response format.
	runner.SendMessage("channel", "?learn call ?response", MsgHelpLearn)
	runner.SendMessage("channel", "?learn call !response", MsgHelpLearn)
	// Valid learns.
	runner.SendLearnMessage("channel", "?learn call response", NewLearn("call", "response"))
	runner.SendLearnMessage("channel", "?learn call2 multi word response", NewLearn("call2", "multi word response"))
	runner.SendLearnMessage("channel", "?learn call3 multi\nline\nresponse\n", NewLearn("call3", "multi\nline\nresponse\n"))
	runner.SendLearnMessage("channel", "?learn call4 \\/leave", NewLearn("call4", "\\/leave"))
	runner.SendLearnMessage("channel", "?learn bearshrug ʅʕ•ᴥ•ʔʃ", NewLearn("bearshrug", "ʅʕ•ᴥ•ʔʃ"))
	runner.SendLearnMessage("channel", "?learn emoji ⛄⛄⛄⛄", NewLearn("emoji", "⛄⛄⛄⛄")) // Emoji is "snowman without snow", in case this isn't showing up in your editor.
	runner.SendLearnMessage("channel", "?learn args1 hello $1", NewLearn("args1", "hello $1"))
	runner.SendLearnMessage("channel", "?learn args2 $1", NewLearn("args2", "$1"))
	runner.SendLearnMessage("channel", "?learn args3 $1 $1", NewLearn("args3", "$1 $1"))
	// Cannot overwrite a learn.
	runner.SendMessage("channel", "?learn call response", fmt.Sprintf(MsgLearnFail, "call"))
	// List should now include learns.
	runner.SendListMessage("channel")
	// Extra whitespace test.
	runner.SendLearnMessage("channel", "?learn  spaceBeforeCall response", NewLearn("spaceBeforeCall", "response"))
	runner.SendLearnMessage("channel", "?learn spaceBeforeResponse  response", NewLearn("spaceBeforeResponse", "response"))
	runner.SendLearnMessage("channel", "?learn spaceInResponse response  two  spaces", NewLearn("spaceInResponse", "response  two  spaces"))

	// Test learned commands.
	runner.SendMessage("channel", "?call", "response")
	runner.SendMessage("channel", "?call2", "multi word response")
	runner.SendMessage("channel", "?call3", "multi\nline\nresponse\n")
	runner.SendMessage("channel", "?call4", "\\/leave")
	runner.SendMessage("channel", "?bearshrug", "ʅʕ•ᴥ•ʔʃ")
	runner.SendMessage("channel", "?emoji", "⛄⛄⛄⛄")
	runner.SendMessage("channel", "?args1 world", "hello world")
	runner.SendMessage("channel", "?args2 world", "world")
	runner.SendMessage("channel", "?args3 world", "world $1")
	runner.SendMessage("channel", "?args3     leadingspaces", "    leadingspaces $1")
	runner.SendMessage("channel", "?args1", MsgCustomNeedsArgs)
	runner.SendMessage("channel", "?spaceBeforeCall", "response")
	runner.SendMessage("channel", "?spaceBeforeResponse", "response")
	runner.SendMessage("channel", "?spaceInResponse", "response  two  spaces")
	// Fallback commands aren't triggered unless they lead a message.
	runner.SendMessageWithoutResponse("channel", " ?call")
	runner.SendMessageWithoutResponse("channel", "i just met you, and this is lazy, but here's my number, ?call me maybe")
	runner.SendMessageWithoutResponse("channel", "\n?call")
	// List should still have the messages.
	runner.SendListMessage("channel")

	// Test unlearn.
	// Wrong format.
	runner.SendMessage("channel", "?unlearn", MsgHelpUnlearn)
	runner.SendMessage("channel", "?unlearn ", MsgHelpUnlearn)
	// Can't unlearn in a private channel
	runner.SendMessage("literally anything else", "?unlearn call", MsgUnlearnMustBePublic)
	// Can't unlearn builtin commands.
	runner.SendMessage("channel", "?unlearn help", fmt.Sprintf(MsgUnlearnFail, "help"))
	runner.SendMessage("channel", "?unlearn learn", fmt.Sprintf(MsgUnlearnFail, "learn"))
	runner.SendMessage("channel", "?unlearn list", fmt.Sprintf(MsgUnlearnFail, "list"))
	runner.SendMessage("channel", "?unlearn unlearn", fmt.Sprintf(MsgUnlearnFail, "unlearn"))
	runner.SendMessage("channel", "?unlearn ?help", MsgHelpUnlearn)
	runner.SendMessage("channel", "?unlearn ?learn", MsgHelpUnlearn)
	runner.SendMessage("channel", "?unlearn ?list", MsgHelpUnlearn)
	runner.SendMessage("channel", "?unlearn ?unlearn", MsgHelpUnlearn)
	// Unrecognized command.
	runner.SendMessage("channel", "?unlearn  bears", fmt.Sprintf(MsgUnlearnFail, "bears"))
	runner.SendMessage("channel", "?unlearn somethingIdon'tknow", fmt.Sprintf(MsgUnlearnFail, "somethingIdon'tknow"))
	// Valid unlearn.
	runner.SendUnlearnMessage("channel", "?unlearn call", "call")
	runner.SendMessageWithoutResponse("channel", "?call")
	// List should work after the unlearn.
	runner.SendListMessage("channel")
	// Can then relearn.
	runner.SendLearnMessage("channel", "?learn call another response", NewLearn("call", "another response"))
	runner.SendMessage("channel", "?call", "another response")
	// List should work after the relearn.
	runner.SendListMessage("channel")
	// Unlearn with 2 spaces.
	runner.SendUnlearnMessage("channel", "?unlearn  call", "call")
	runner.SendMessageWithoutResponse("channel", "?call")

	// Unrecognized help commands.
	runner.SendMessage("channel", "?help", MsgDefaultHelp)
	runner.SendMessage("channel", "?help abunchofgibberish", MsgDefaultHelp)
	runner.SendMessage("channel", "?help ??help", MsgDefaultHelp)
	// All recognized help commands.
	runner.SendMessage("channel", "?help help", MsgHelpHelp)
	runner.SendMessage("channel", "?help learn", MsgHelpLearn)
	runner.SendMessage("channel", "?help list", MsgHelpList)
	runner.SendMessage("channel", "?help unlearn", MsgHelpUnlearn)
	runner.SendMessage("channel", "?help ?help", MsgHelpHelp)
	runner.SendMessage("channel", "?help ?learn", MsgHelpLearn)
	runner.SendMessage("channel", "?help ?list", MsgHelpList)
	runner.SendMessage("channel", "?help ?unlearn", MsgHelpUnlearn)
	runner.SendMessage("channel", "?help  help", MsgHelpHelp)
	// Help with custom commands.
	runner.SendLearnMessage("channel", "?learn help-noarg response", NewLearn("help-noarg", "response"))
	runner.SendLearnMessage("channel", "?learn help-arg response $1", NewLearn("help-arg", "response $1"))
	runner.SendMessage("channel", "?help help-noarg", "?help-noarg")
	runner.SendMessage("channel", "?help help-arg", "?help-arg <args>")
	runner.SendUnlearnMessage("channel", "?unlearn help-noarg", "help-noarg")
	runner.SendUnlearnMessage("channel", "?unlearn help-arg", "help-arg")
	runner.SendMessage("channel", "?help help-noarg", MsgDefaultHelp)
	runner.SendMessage("channel", "?help help-arg", MsgDefaultHelp)
}

func assertNumCommands(t *testing.T, customMap StringMap, count int) {
	if all, _ := customMap.GetAll(); len(all) != count {
		t.Errorf(fmt.Sprintf("Should have %v commands", count))
	}
}

func assertNumGists(t *testing.T, gist *util.InMemoryGist, count int) {
	if len(gist.Messages) != count {
		t.Errorf(fmt.Sprintf("Should have %v gists", count))
	}
}

func assertNumDiscordMessages(t *testing.T, discordSession *util.InMemoryDiscordSession, count int) {
	if len(discordSession.Messages) != count {
		t.Errorf(fmt.Sprintf("Should have %v discord messages", count))
	}
}

func sendMessage(t *testing.T, discordSession DiscordSession, handler func(DiscordSession, *discordgo.MessageCreate), channel, message string) {
	author := &discordgo.User{
		ID:            "author id",
		Email:         "email@example.com",
		Username:      "username",
		Avatar:        "avatar",
		Discriminator: "discriminator",
		Token:         "token",
		Verified:      true,
		MFAEnabled:    false,
		Bot:           false,
	}

	messageCreate := &discordgo.MessageCreate{
		&discordgo.Message{
			ID:              "messageID",
			ChannelID:       channel,
			Content:         message,
			Timestamp:       "timestamp",
			EditedTimestamp: "edited timestamp",
			MentionRoles:    []string{},
			Tts:             false,
			MentionEveryone: false,
			Author:          author,
			Attachments:     []*discordgo.MessageAttachment{},
			Embeds:          []*discordgo.MessageEmbed{},
			Mentions:        []*discordgo.User{},
			Reactions:       []*discordgo.MessageReactions{},
		},
	}
	handler(discordSession, messageCreate)
}

func assertNewMessages(t *testing.T, discordSession *util.InMemoryDiscordSession, newMessages []*util.Message) {
	if len(discordSession.Messages) < len(newMessages) {
		t.Errorf(fmt.Sprintf(
			"Needed at least %v messages, had %v", len(newMessages), len(discordSession.Messages)))
		return
	}

	for i := 0; i < len(newMessages); i++ {
		actualMessage := discordSession.Messages[len(discordSession.Messages)-len(newMessages)+i]
		if !reflect.DeepEqual(newMessages[i], actualMessage) {
			t.Errorf("Expected message %v channel %v, got message %v channel %v",
				newMessages[i].Message,
				newMessages[i].Channel,
				actualMessage.Message,
				actualMessage.Channel)
		}
	}
}

func assertCommand(t *testing.T, commandMap *util.InMemoryStringMap, call, response string) {
	if _, err := commandMap.Get(call); err != nil {
		t.Errorf("Response should be present for call " + call)
		return
	}
	if value, _ := commandMap.Get(call); value != response {
		t.Errorf(fmt.Sprintf("Wrong response for %v, expected %v got %v", call, response, value))
	}
}

type Learn struct {
	Call     string
	Response string
}

func NewLearn(call, response string) *Learn {
	return &Learn{
		Call:     call,
		Response: response,
	}
}
