package testutil

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aetimmes/discordgo"
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

	registry := app.InitializeRegistry(customMap, karmaMap, voteMap, gist, &config.Config{RickList: rickList}, utcClock, utcTimer, commandChannel)

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
	assertVote(r.T, r.UTCClock, r.VoteMap, r.ActiveVoteDataMap)

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

// SendVoteMessageAs sends a ?vote message to the bot as the given user
func (r *Runner) SendVoteMessageAs(author *discordgo.User, channel model.Snowflake) {
	r.T.Helper()

	sendMessageAs(author, r.DiscordSession, r.Handler, channel, "?vote a vote has been called")
	r.DiscordMessagesCount++
	r.ActiveVoteDataMap[channel] = newVoteData(channel, author, "a vote has been called", r.UTCClock.Now().Add(vote.VoteDuration))
	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), fmt.Sprintf(vote.MsgBroadcastNewVote, author.Mention(), "a vote has been called"))})
	r.AssertState()
}

// CastBallotAs casts a ballot as the given user
func (r *Runner) CastBallotAs(author *discordgo.User, channel model.Snowflake, inFavor bool) {
	r.T.Helper()

	voteString := "?no"
	expectedMessage := fmt.Sprintf(vote.MsgVotedAgainst, author.Mention())
	toAppend := &(r.ActiveVoteDataMap[channel].VotesAgainst)
	if inFavor {
		voteString = "?yes"
		expectedMessage = fmt.Sprintf(vote.MsgVotedInFavor, author.Mention())
		toAppend = &(r.ActiveVoteDataMap[channel].VotesFor)
	}

	sendMessageAs(author, r.DiscordSession, r.Handler, channel, voteString)

	// Update internal state.
	r.DiscordMessagesCount++
	id, _ := model.ParseSnowflake(author.ID)
	*toAppend = append(*toAppend, id)

	// Reconstruct the status string and assert internal state.
	activeVote := r.ActiveVoteDataMap[channel]
	reconstructedVote := activeVote.Reconstruct()

	assertNewMessages(r.T, r.DiscordSession, []*Message{
		NewMessage(channel.Format(), expectedMessage+"\n"+vote.StatusLine(r.UTCClock, reconstructedVote)),
	})
	r.AssertState()
}

// CastDuplicateBallotAs casts a ballot as the user and expects no state change as a result
func (r *Runner) CastDuplicateBallotAs(author *discordgo.User, channel model.Snowflake, inFavor bool) {
	r.T.Helper()

	voteString := "?no"
	if inFavor {
		voteString = "?yes"
	}

	sendMessageAs(author, r.DiscordSession, r.Handler, channel, voteString)

	// Update internal state.
	r.DiscordMessagesCount++

	assertNewMessages(r.T, r.DiscordSession, []*Message{
		NewMessage(channel.Format(), fmt.Sprintf(vote.MsgAlreadyVoted, author.Mention())),
	})
	r.AssertState()
}

// ExpireVote advances the clock enough that the vote expires, and fires the trigger.
func (r *Runner) ExpireVote(channel model.Snowflake) {
	r.T.Helper()

	v := r.ActiveVoteDataMap[channel]

	// Calculate the time to elapse to expire the given vote.
	toElapse := v.TimestampEnd.Sub(r.UTCClock.Now())

	// The UTC clock and UTC timer need to be advanced together.
	r.UTCClock.Advance(toElapse)

	// Elapse time requires a flush because it can generate the conclude command.
	r.UTCTimer.ElapseTime(toElapse)
	flushChannel(r.DiscordSession, r.Handler, channel)

	// Update internal state.
	r.DiscordMessagesCount++

	// Calculate expected vote outcome.
	voteOutcome := model.VoteOutcomeNotEnough
	if len(v.VotesFor)+len(v.VotesAgainst) >= 5 {
		voteOutcome = model.VoteOutcomeFailed
		if len(v.VotesFor) > len(v.VotesAgainst) {
			voteOutcome = model.VoteOutcomePassed
		}
	}

	reconstructedVote := v.Reconstruct()
	reconstructedVote.VoteOutcome = voteOutcome
	statusLine := vote.CompletedStatusLine(reconstructedVote)

	expectedMessage := fmt.Sprintf(vote.MsgVoteConcluded, v.Author.Mention()) +
		"\n" +
		statusLine

	r.ActiveVoteDataMap[channel] = nil

	assertNewMessages(r.T, r.DiscordSession,
		[]*Message{NewMessage(channel.Format(), expectedMessage)})
	r.AssertState()
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
		buffer.WriteString(" - ?f1: ")
		buffer.WriteString(vote.MsgHelpBallotInFavor)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?f2: ")
		buffer.WriteString(vote.MsgHelpBallotAgainst)
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
		buffer.WriteString(" - ?no: ")
		buffer.WriteString(vote.MsgHelpBallotAgainst)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?ricklist: ")
		buffer.WriteString(moderation.MsgHelpRickListInfo)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?unlearn: ")
		buffer.WriteString(learn.MsgHelpUnlearn)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?vote: ")
		buffer.WriteString(vote.MsgHelpVote)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?votestatus: ")
		buffer.WriteString(vote.MsgHelpStatus)
		buffer.WriteString("\n")
		buffer.WriteString(" - ?yes: ")
		buffer.WriteString(vote.MsgHelpBallotInFavor)
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

// SendVoteStatusMessage sends a vote status message to the bot
func (r *Runner) SendVoteStatusMessage(channel model.Snowflake) {
	r.T.Helper()

	sendMessage(r.DiscordSession, r.Handler, channel, "?votestatus")
	r.DiscordMessagesCount++

	activeVote, _ := r.ActiveVoteDataMap[channel]

	if activeVote == nil {
		assertNewMessages(r.T, r.DiscordSession, []*Message{NewMessage(channel.Format(), vote.MsgNoActiveVote)})
	} else {
		// Calculate the expected status messages.
		forMessage := vote.MsgOneVoteFor
		if len(activeVote.VotesFor) != 1 {
			forMessage = fmt.Sprintf(vote.MsgVotesFor, len(activeVote.VotesFor))
		}
		againstMessage := vote.MsgOneVoteAgainst
		if len(activeVote.VotesAgainst) != 1 {
			againstMessage = fmt.Sprintf(vote.MsgVotesAgainst, len(r.ActiveVoteDataMap[channel].VotesAgainst))
		}
		statusMessage := vote.MsgStatusVotesNeeded
		if len(activeVote.VotesAgainst)+len(activeVote.VotesFor) >= 5 {
			if len(activeVote.VotesFor) > len(activeVote.VotesAgainst) {
				statusMessage = vote.MsgStatusVotePassing
			} else {
				statusMessage = vote.MsgStatusVoteFailing
			}
		}

		// The time remaining is independently tested, so just assert its presence.
		timeMessage := vote.TimeString(r.UTCClock, activeVote.TimestampEnd)

		// Build the expected string and assert that it's in the message buffer.
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf(vote.MsgVoteOwner, activeVote.Author.Username))
		buffer.WriteString(activeVote.Message)
		buffer.WriteString("\n")
		buffer.WriteString(vote.MsgSpacer)
		buffer.WriteString("\n")
		buffer.WriteString(statusMessage)
		buffer.WriteString(". ")
		buffer.WriteString(forMessage)
		buffer.WriteString(", ")
		buffer.WriteString(againstMessage)
		buffer.WriteString(". ")
		buffer.WriteString(timeMessage)
		assertNewMessages(r.T, r.DiscordSession, []*Message{NewMessage(channel.Format(), buffer.String())})
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

// Reconstruct creates a vote.Vote out of the local storage vote. This isn't
// meant to be a complete reconstruction, but rather all the info necessary for
// testing (mostly reproducing status lines).
func (v *VoteData) Reconstruct() *model.Vote {
	parsedSnowflake, _ := model.ParseSnowflake(v.Author.ID)
	return model.NewVote(
		0, /* voteID */
		v.Channel,
		parsedSnowflake,
		v.Message,
		time.Time{},
		v.TimestampEnd,
		v.VotesFor,
		v.VotesAgainst,
		model.VoteOutcomeNotDone)
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

func assertVote(t *testing.T, utcClock model.UTCClock, voteMap *stringmap.InMemoryStringMap, activeVoteMap map[model.Snowflake]*VoteData) {
	t.Helper()

	modelHelper := vote.NewModelHelper(voteMap, utcClock)
	for channel, vote := range activeVoteMap {
		ok, _ := modelHelper.IsVoteActive(channel)
		if vote != nil && !ok {
			t.Errorf("Expected a vote to be active, but was not")
		}
		if vote == nil && ok {
			t.Errorf("Expected a vote to not be active, but one was")
		}
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
	messageCreate := &discordgo.MessageCreate{ // nolint
		&discordgo.Message{
			ID:              "messageID",
			ChannelID:       channel.Format(),
			Content:         message,
			Timestamp:       "timestamp",
			EditedTimestamp: "edited timestamp",
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
		&discordgo.Message{
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
