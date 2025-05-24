package testutil

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/config"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/feature/factsphere"
	"github.com/jakevoytko/crbot/feature/help"
	"github.com/jakevoytko/crbot/feature/karma"
	"github.com/jakevoytko/crbot/feature/karmalist"
	"github.com/jakevoytko/crbot/feature/learn"
	"github.com/jakevoytko/crbot/feature/list"
	"github.com/jakevoytko/crbot/feature/moderation"
	"github.com/jakevoytko/crbot/feature/vote"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// IDs used for testing
const (
	MainChannelID   = model.Snowflake(8675309)
	SecondChannelID = model.Snowflake(9000000)
	DirectMessageID = model.Snowflake(1)
)

// Runner is a helper that executes messages incrementally, and asserts that
// the global state is always what is expected.
type Runner struct {
	// Test object
	T *testing.T

	// State
	GistsCount           int
	DiscordMessagesCount int
	LearnDataMap         map[string]*LearnData
	ActiveVoteDataMap    map[model.Snowflake]*VoteData // channel->vote. May be nil

	// Fakes
	CustomMap      *stringmap.InMemoryStringMap
	KarmaMap       *stringmap.InMemoryStringMap
	VoteMap        *stringmap.InMemoryStringMap
	Gist           *InMemoryGist
	DiscordSession *InMemoryDiscordSession
	UTCClock       *FakeUTCClock
	UTCTimer       *FakeUTCTimer

	// Real objects
	FeatureRegistry *feature.Registry

	// Controllers under test
	Handler func(api.DiscordSession, *discordgo.MessageCreate)
}

// NewRunner works as advertised
func NewRunner(t *testing.T) *Runner {
	// Initialize fakes.
	customMap := stringmap.NewInMemoryStringMap()
	karmaMap := stringmap.NewInMemoryStringMap()
	voteMap := stringmap.NewInMemoryStringMap()
	gist := NewInMemoryGist()
	discordSession := NewInMemoryDiscordSession()
	discordSession.SetChannel(&discordgo.Channel{
		ID:   MainChannelID.Format(),
		Type: discordgo.ChannelTypeGuildText,
	})
	discordSession.SetChannel(&discordgo.Channel{
		ID:   SecondChannelID.Format(),
		Type: discordgo.ChannelTypeGuildText,
	})
	discordSession.SetChannel(&discordgo.Channel{
		ID:   DirectMessageID.Format(),
		Type: discordgo.ChannelTypeDM,
	})

	rickList := []model.Snowflake{model.Snowflake(2)}

	utcClock := NewFakeUTCClock()

	// 0-length channel. Each time it consumes/processes a command, it issues a
	// no-op as a hack to make sure that the first one has been processed.
	commandChannel := make(chan *model.Command)

	utcTimer := NewFakeUTCTimer()

	registry := app.InitializeRegistry(customMap, karmaMap, gist, &config.Config{RickList: rickList})

	go app.HandleCommands(registry, discordSession, commandChannel)

	return &Runner{
		T:                    t,
		LearnDataMap:         map[string]*LearnData{},
		ActiveVoteDataMap:    map[model.Snowflake]*VoteData{},
		GistsCount:           0,
		DiscordMessagesCount: 0,
		CustomMap:            customMap,
		KarmaMap:             karmaMap,
		VoteMap:              voteMap,
		Gist:                 gist,
		DiscordSession:       discordSession,
		UTCClock:             utcClock,
		UTCTimer:             utcTimer,
		FeatureRegistry:      registry,
		Handler:              app.GetHandleMessage(customMap, registry, commandChannel),
	}
}

// AssertState asserts all of the features against the test runner's state
func (r *Runner) AssertState() {
	r.T.Helper()

	// Assert counts.
	assertNumCommands(r.T, r.CustomMap, len(r.LearnDataMap))
	assertNumGists(r.T, r.Gist, r.GistsCount)
	assertNumDiscordMessages(r.T, r.DiscordSession, r.DiscordMessagesCount)

	// Assert command map state.
	for _, learn := range r.LearnDataMap {
		assertCommand(r.T, r.CustomMap, learn.Call, learn.Response)
	}
}

// SendMessage sends the message to the bot as the standard test user
func (r *Runner) SendMessage(channel model.Snowflake, message, expectedResponse string) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), expectedResponse)})
	r.AssertState()
}

// SendMessageAs sends a message to the bot as the given user
func (r *Runner) SendMessageAs(author *discordgo.User, channel model.Snowflake, message, expectedResponse string) {
	r.T.Helper()

	sendMessageAs(author, r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), expectedResponse)})
	r.AssertState()
}

// SendMessageIgnoringResponse sends a message to the bot without checking the output
func (r *Runner) SendMessageIgnoringResponse(channel model.Snowflake, message string) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	r.AssertState()
}

// SendLearnMessage sends a learn message to the bot
func (r *Runner) SendLearnMessage(channel model.Snowflake, message string, learnData *LearnData) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	r.LearnDataMap[learnData.Call] = learnData
	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), fmt.Sprintf(learn.MsgLearnSuccess, learnData.Call))})
	r.AssertState()
	r.SendListMessage(channel)
}

// SendLearnMessageAs sends a ?learn message as the given user
func (r *Runner) SendLearnMessageAs(author *discordgo.User, channel model.Snowflake, message string, learnData *LearnData) {
	r.T.Helper()

	sendMessageAs(author, r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	r.LearnDataMap[learnData.Call] = learnData
	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), fmt.Sprintf(learn.MsgLearnSuccess, learnData.Call))})
	r.AssertState()
	r.SendListMessage(channel)
}

// SendUnlearnMessage sends an ?unlearn message
func (r *Runner) SendUnlearnMessage(channel model.Snowflake, message string, call string) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, message)
	r.DiscordMessagesCount++
	delete(r.LearnDataMap, call)
	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), fmt.Sprintf(learn.MsgUnlearnSuccess, call))})
	r.AssertState()
	r.SendListMessage(channel)
}

// SendMessageWithoutResponse sends a message without a response
func (r *Runner) SendMessageWithoutResponse(channel model.Snowflake, message string) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, message)
	assertNewMessages(r.T, r.DiscordSession, []*Message{})
	r.AssertState()
}

// SendListMessage sends a ?list message to the bot
func (r *Runner) SendListMessage(channel model.Snowflake) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, "?list")
	r.DiscordMessagesCount++
	r.GistsCount++
	assertNewMessages(r.T, r.DiscordSession, []*Message{NewMessage(channel.Format(), "The list of commands is here: https://www.example.com/success")})

	// Assert gist state. Cannot be in AssertState because this would fail at the
	// next learn or unlearn.
	// TODO(jake): Remove duplication between this and listfeature. Maybe just assert number of lines?
	if r.GistsCount > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("List of builtins:")
		buffer.WriteString("\n")
		buffer.WriteString(" - ?++: ")
		buffer.WriteString(karma.MsgHelpKarmaIncrement)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?--: ")
		buffer.WriteString(karma.MsgHelpKarmaDecrement)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?factsphere: ")
		buffer.WriteString(factsphere.MsgHelpFactSphere)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?help: ")
		buffer.WriteString(help.MsgHelpHelp)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?karmalist: ")
		buffer.WriteString(karmalist.MsgHelpKarmaList)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?learn: ")
		buffer.WriteString(learn.MsgHelpLearn)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?list: ")
		buffer.WriteString(list.MsgHelpList)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?ricklist: ")
		buffer.WriteString(moderation.MsgHelpRickListInfo)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?unlearn: ")
		buffer.WriteString(learn.MsgHelpUnlearn)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?vote: ")
		buffer.WriteString(vote.MsgHelpStatus)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?votestatus: ")
		buffer.WriteString(vote.MsgHelpStatus)
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

// SendKarmaListMessage senda a ?karmalist message to the bot
func (r *Runner) SendKarmaListMessage(channel model.Snowflake) {
	r.T.Helper()
	sendMessage(r.DiscordSession, r.Handler, channel, "?karmalist")
	r.DiscordMessagesCount++
	r.GistsCount++
	assertNewMessages(r.T, r.DiscordSession, []*Message{NewMessage(channel.Format(), "The list of karma is here: https://www.example.com/success")})
	if r.GistsCount > 0 {
		karmaRunner := karmalist.NewModelHelper(r.KarmaMap)
		generated := karmaRunner.GenerateList()
		actual := r.Gist.Messages[len(r.Gist.Messages)-1]
		if generated != actual {
			r.T.Fatalf(fmt.Sprintf("Gist failure, got `%v` expected `%v`", actual, generated))
		}
	}

	r.AssertState()
}

// AddUser adds a user in the test session
func (r *Runner) AddUser(user *discordgo.User) {
	r.DiscordSession.Users[user.ID] = user
}

// LearnData represents the information needed to reconstruct a Learn for testing.
type LearnData struct {
	Call     string
	Response string
}

// NewLearnData stores a duplicate representation of the expected learn data for testing
func NewLearnData(call, response string) *LearnData {
	return &LearnData{
		Call:     call,
		Response: response,
	}
}

// VoteData contains enough information to reconstruct the status message.
type VoteData struct {
	Channel      model.Snowflake
	Author       *discordgo.User
	Message      string
	VotesFor     []model.Snowflake
	VotesAgainst []model.Snowflake
	TimestampEnd time.Time
}

func newVoteData(channel model.Snowflake, author *discordgo.User, message string, timestampEnd time.Time) *VoteData {
	return &VoteData{
		Channel:      channel,
		Author:       author,
		Message:      message,
		VotesFor:     []model.Snowflake{},
		VotesAgainst: []model.Snowflake{},
		TimestampEnd: timestampEnd,
	}
}

func assertNumCommands(t *testing.T, customMap *stringmap.InMemoryStringMap, count int) {
	t.Helper()

	if all, _ := customMap.GetAll(); len(all) != count {
		t.Errorf(fmt.Sprintf("Should have %v commands", count))
	}
}

func assertNumGists(t *testing.T, gist *InMemoryGist, count int) {
	t.Helper()

	if len(gist.Messages) != count {
		t.Errorf(fmt.Sprintf("Should have %v gists", count))
	}
}

func assertNumDiscordMessages(t *testing.T, discordSession *InMemoryDiscordSession, count int) {
	t.Helper()

	if len(discordSession.Messages) != count {
		t.Errorf(fmt.Sprintf("Should have %v discord messages", count))
	}
}

func sendMessage(discordSession api.DiscordSession, handler func(api.DiscordSession, *discordgo.MessageCreate), channel model.Snowflake, message string) {
	author := &discordgo.User{
		ID:            "1",
		Email:         "email@example.com",
		Username:      "username",
		Avatar:        "avatar",
		Discriminator: "discriminator",
		Token:         "token",
		Verified:      true,
		MFAEnabled:    false,
		Bot:           false,
	}

	sendMessageAs(author, discordSession, handler, channel, message)
}

func sendMessageAs(author *discordgo.User, discordSession api.DiscordSession, handler func(api.DiscordSession, *discordgo.MessageCreate), channel model.Snowflake, message string) {
	editedTimestamp := time.Now()
	messageCreate := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:              "messageID",
			ChannelID:       channel.Format(),
			Content:         message,
			Timestamp:       time.Now().Add(-time.Hour),
			EditedTimestamp: &editedTimestamp,
			MentionRoles:    []string{},
			TTS:             false,
			MentionEveryone: false,
			Author:          author,
			Attachments:     []*discordgo.MessageAttachment{},
			Embeds:          []*discordgo.MessageEmbed{},
			Mentions:        []*discordgo.User{},
			Reactions:       []*discordgo.MessageReactions{},
		},
	}
	handler(discordSession, messageCreate)
	flushChannel(discordSession, handler, channel)
}

func flushChannel(discordSession api.DiscordSession, handler func(api.DiscordSession, *discordgo.MessageCreate), channel model.Snowflake) {
	user := NewUser("fake", model.Snowflake(0), false /* isBot */)

	// A no-op command that flushes out the 0 length buffer so assertions are
	// correct. Otherwise, processing would happen asynchronously, so it'd be
	// impossible to assert that the program had behaved correctly.
	noOp := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:    user,
			ChannelID: channel.Format(),
			Content:   "",
		},
	}
	handler(discordSession, noOp)
}

func assertNewMessages(t *testing.T, discordSession *InMemoryDiscordSession, newMessages []*Message) {
	t.Helper()

	if len(discordSession.Messages) < len(newMessages) {
		t.Errorf(fmt.Sprintf(
			"Needed at least %v messages, had %v", len(newMessages), len(discordSession.Messages)))
		return
	}

	for i := 0; i < len(newMessages); i++ {
		actualMessage := discordSession.Messages[len(discordSession.Messages)-len(newMessages)+i]
		if !reflect.DeepEqual(newMessages[i], actualMessage) {
			t.Errorf("Expected message \n '%v' \n on channel '%v', got message \n '%v' \n on channel '%v'",
				newMessages[i].Message,
				newMessages[i].Channel,
				actualMessage.Message,
				actualMessage.Channel)
		}
	}
}

func assertCommand(t *testing.T, commandMap *stringmap.InMemoryStringMap, call, response string) {
	t.Helper()

	if _, err := commandMap.Get(call); err != nil {
		t.Errorf("Response should be present for call " + call)
		return
	}
	if value, _ := commandMap.Get(call); value != response {
		t.Errorf(fmt.Sprintf("Wrong response for %v, expected %v got %v", call, response, value))
	}
}

// NewUser creates a new user for testing
func NewUser(name string, id model.Snowflake, bot bool) *discordgo.User {
	idStr := id.Format()
	return &discordgo.User{
		ID:            idStr,
		Email:         "email@example.com",
		Username:      name,
		Avatar:        "avatar",
		Discriminator: idStr,
		Token:         "token",
		Verified:      true,
		MFAEnabled:    false,
		Bot:           bot,
	}
}
