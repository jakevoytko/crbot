package main

import (
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/feature/help"
	"github.com/jakevoytko/crbot/feature/learn"
	"github.com/jakevoytko/crbot/feature/list"
	"github.com/jakevoytko/crbot/feature/moderation"
	"github.com/jakevoytko/crbot/feature/vote"
	"github.com/jakevoytko/crbot/testutil"
)

func TestLearn_NoResponse(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Commands that should never return a response.
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "?")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "!")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, ".")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "!help")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "help")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, ".help")

	// Test ?list. ?list tests will be interspersed through the learn examples
	// below, since learn and unlearn interact with it.
	runner.SendListMessage(testutil.MainChannelID)
}

func TestLearn_WrongFormat(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Basic learn responses.
	// Wrong call format
	runner.SendMessage(testutil.MainChannelID, "?learn", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn test", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn ?call response", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn !call response", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn /call response", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn ", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn multi\nline\ncall response", learn.MsgHelpLearn)
	// Wrong response format.
	runner.SendMessage(testutil.MainChannelID, "?learn call ?response", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?learn call !response", learn.MsgHelpLearn)
}

func TestIntegration(t *testing.T) {
	runner := testutil.NewRunner(t)
	// Test ?list. ?list tests will be interspersed through the learn examples
	// below, since learn and unlearn interact with it.
	runner.SendListMessage(testutil.MainChannelID)

	// Valid learns.
	runner.SendLearnMessage(testutil.MainChannelID, "?learn call response", testutil.NewLearnData("call", "response"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn call2 multi word response", testutil.NewLearnData("call2", "multi word response"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn call3 multi\nline\nresponse\n", testutil.NewLearnData("call3", "multi\nline\nresponse\n"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn call4 \\/leave", testutil.NewLearnData("call4", "\\/leave"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn bearshrug ʅʕ•ᴥ•ʔʃ", testutil.NewLearnData("bearshrug", "ʅʕ•ᴥ•ʔʃ"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn emoji ⛄⛄⛄⛄", testutil.NewLearnData("emoji", "⛄⛄⛄⛄")) // Emoji is "snowman without snow", in case this isn't showing up in your editor.
	runner.SendLearnMessage(testutil.MainChannelID, "?learn args1 hello $1", testutil.NewLearnData("args1", "hello $1"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn args2 $1", testutil.NewLearnData("args2", "$1"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn args3 $1 $1", testutil.NewLearnData("args3", "$1 $1"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn args4 $1 $1 $1 $1 $1", testutil.NewLearnData("args4", "$1 $1 $1 $1 $1"))
	// Cannot overwrite a learn.
	runner.SendMessage(testutil.MainChannelID, "?learn call response", fmt.Sprintf(learn.MsgLearnFail, "call"))
	// List should now include learns.
	runner.SendListMessage(testutil.MainChannelID)
	// Extra whitespace test.
	runner.SendLearnMessage(testutil.MainChannelID, "?learn  spaceBeforeCall response", testutil.NewLearnData("spaceBeforeCall", "response"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn spaceBeforeResponse  response", testutil.NewLearnData("spaceBeforeResponse", "response"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn spaceInResponse response  two  spaces", testutil.NewLearnData("spaceInResponse", "response  two  spaces"))

	// Test learned commands.
	runner.SendMessage(testutil.MainChannelID, "?call", "response")
	runner.SendMessage(testutil.MainChannelID, "?call2", "multi word response")
	runner.SendMessage(testutil.MainChannelID, "?call3", "multi\nline\nresponse\n")
	runner.SendMessage(testutil.MainChannelID, "?call4", "\\/leave")
	runner.SendMessage(testutil.MainChannelID, "?bearshrug", "ʅʕ•ᴥ•ʔʃ")
	runner.SendMessage(testutil.MainChannelID, "?emoji", "⛄⛄⛄⛄")
	runner.SendMessage(testutil.MainChannelID, "?args1 world", "hello world")
	runner.SendMessage(testutil.MainChannelID, "?args2 world", "world")
	runner.SendMessage(testutil.MainChannelID, "?args3 world", "world world")
	runner.SendMessage(testutil.MainChannelID, "?args3     leadingspaces", "    leadingspaces     leadingspaces")
	runner.SendMessage(testutil.MainChannelID, "?args4 world", "world world world world $1")
	runner.SendMessage(testutil.MainChannelID, "?args4     leadingspaces", "    leadingspaces     leadingspaces     leadingspaces     leadingspaces $1")

	runner.SendMessage(testutil.MainChannelID, "?args1", learn.MsgCustomNeedsArgs)
	runner.SendMessage(testutil.MainChannelID, "?spaceBeforeCall", "response")
	runner.SendMessage(testutil.MainChannelID, "?spaceBeforeResponse", "response")
	runner.SendMessage(testutil.MainChannelID, "?spaceInResponse", "response  two  spaces")
	// Fallback commands aren't triggered unless they lead a message.
	runner.SendMessageWithoutResponse(testutil.MainChannelID, " ?call")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "i just met you, and this is lazy, but here's my number, ?call me maybe")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "\n?call")
	// List should still have the messages.
	runner.SendListMessage(testutil.MainChannelID)

	// Test unlearn.
	// Wrong format.
	runner.SendMessage(testutil.MainChannelID, "?unlearn", learn.MsgHelpUnlearn)
	runner.SendMessage(testutil.MainChannelID, "?unlearn ", learn.MsgHelpUnlearn)
	// Can't unlearn in a private channel
	runner.SendMessage(testutil.DirectMessageID, "?unlearn call", learn.MsgUnlearnMustBePublic)
	// Can't unlearn builtin commands.
	runner.SendMessage(testutil.MainChannelID, "?unlearn help", fmt.Sprintf(learn.MsgUnlearnFail, "help"))
	runner.SendMessage(testutil.MainChannelID, "?unlearn learn", fmt.Sprintf(learn.MsgUnlearnFail, "learn"))
	runner.SendMessage(testutil.MainChannelID, "?unlearn list", fmt.Sprintf(learn.MsgUnlearnFail, "list"))
	runner.SendMessage(testutil.MainChannelID, "?unlearn unlearn", fmt.Sprintf(learn.MsgUnlearnFail, "unlearn"))
	runner.SendMessage(testutil.MainChannelID, "?unlearn ?help", learn.MsgHelpUnlearn)
	runner.SendMessage(testutil.MainChannelID, "?unlearn ?learn", learn.MsgHelpUnlearn)
	runner.SendMessage(testutil.MainChannelID, "?unlearn ?list", learn.MsgHelpUnlearn)
	runner.SendMessage(testutil.MainChannelID, "?unlearn ?unlearn", learn.MsgHelpUnlearn)
	// Unrecognized command.
	runner.SendMessage(testutil.MainChannelID, "?unlearn  bears", fmt.Sprintf(learn.MsgUnlearnFail, "bears"))
	runner.SendMessage(testutil.MainChannelID, "?unlearn somethingIdon'tknow", fmt.Sprintf(learn.MsgUnlearnFail, "somethingIdon'tknow"))
	// Valid unlearn.
	runner.SendUnlearnMessage(testutil.MainChannelID, "?unlearn call", "call")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "?call")
	// List should work after the unlearn.
	runner.SendListMessage(testutil.MainChannelID)
	// Can then relearn.
	runner.SendLearnMessage(testutil.MainChannelID, "?learn call another response", testutil.NewLearnData("call", "another response"))
	runner.SendMessage(testutil.MainChannelID, "?call", "another response")
	// List should work after the relearn.
	runner.SendListMessage(testutil.MainChannelID)
	// Unlearn with 2 spaces.
	runner.SendUnlearnMessage(testutil.MainChannelID, "?unlearn  call", "call")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "?call")

	// Unrecognized help commands.
	runner.SendMessage(testutil.MainChannelID, "?help", help.MsgDefaultHelp)
	runner.SendMessage(testutil.MainChannelID, "?help abunchofgibberish", help.MsgDefaultHelp)
	runner.SendMessage(testutil.MainChannelID, "?help ??help", help.MsgDefaultHelp)
	// All recognized help commands.
	runner.SendMessage(testutil.MainChannelID, "?help help", help.MsgHelpHelp)
	runner.SendMessage(testutil.MainChannelID, "?help learn", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?help list", list.MsgHelpList)
	runner.SendMessage(testutil.MainChannelID, "?help unlearn", learn.MsgHelpUnlearn)
	runner.SendMessage(testutil.MainChannelID, "?help ?help", help.MsgHelpHelp)
	runner.SendMessage(testutil.MainChannelID, "?help ?learn", learn.MsgHelpLearn)
	runner.SendMessage(testutil.MainChannelID, "?help ?list", list.MsgHelpList)
	runner.SendMessage(testutil.MainChannelID, "?help ?unlearn", learn.MsgHelpUnlearn)
	runner.SendMessage(testutil.MainChannelID, "?help  help", help.MsgHelpHelp)
	// Help with custom commands.
	runner.SendLearnMessage(testutil.MainChannelID, "?learn help-noarg response", testutil.NewLearnData("help-noarg", "response"))
	runner.SendLearnMessage(testutil.MainChannelID, "?learn help-arg response $1", testutil.NewLearnData("help-arg", "response $1"))
	runner.SendMessage(testutil.MainChannelID, "?help help-noarg", "?help-noarg")
	runner.SendMessage(testutil.MainChannelID, "?help help-arg", "?help-arg <args>")
	runner.SendUnlearnMessage(testutil.MainChannelID, "?unlearn help-noarg", "help-noarg")
	runner.SendUnlearnMessage(testutil.MainChannelID, "?unlearn help-arg", "help-arg")
	runner.SendMessage(testutil.MainChannelID, "?help help-noarg", help.MsgDefaultHelp)
	runner.SendMessage(testutil.MainChannelID, "?help help-arg", help.MsgDefaultHelp)

	// Moderation
	rickListedUser := &discordgo.User{
		ID:            "2",
		Email:         "email@example.com",
		Username:      "username",
		Avatar:        "avatar",
		Discriminator: "discriminator",
		Token:         "token",
		Verified:      true,
		MFAEnabled:    false,
		Bot:           false,
	}
	runner.SendMessageAs(rickListedUser, testutil.MainChannelID, "?help help-arg", help.MsgDefaultHelp)
	runner.SendMessageAs(rickListedUser, testutil.DirectMessageID, "?help help-arg", moderation.MsgRickList)
	runner.SendLearnMessageAs(rickListedUser, testutil.DirectMessageID, "?learn rick list", testutil.NewLearnData("rick", "list"))
}

func TestVote(t *testing.T) {
	runner := testutil.NewRunner(t)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Calls vote with no args, and then actually starts a vote.
	author := testutil.NewUser("author", 0 /* id */, false /* bot */)
	runner.AddUser(author)
	runner.SendMessageAs(author, testutil.MainChannelID, "?vote", vote.MsgHelpVote)
	runner.SendVoteMessageAs(author, testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Assert that a second vote can't be started.
	runner.SendMessageAs(author, testutil.MainChannelID, "?vote another vote", vote.MsgActiveVote)

	// Time the vote out.
	runner.ExpireVote(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// A second vote can be started once it is expired.
	runner.SendVoteMessageAs(author, testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)
}

func TestVote_Pass(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Initialize users.
	users := []*discordgo.User{
		testutil.NewUser("user0", 0 /* id */, false /* bot */),
		testutil.NewUser("user1", 1 /* id */, false /* bot */),
		testutil.NewUser("user2", 2 /* id */, false /* bot */),
		testutil.NewUser("user3", 3 /* id */, false /* bot */),
		testutil.NewUser("user4", 4 /* id */, false /* bot */),
	}
	for _, user := range users {
		runner.AddUser(user)
	}

	// Start the vote.
	runner.SendVoteMessageAs(users[0], testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Cast votes
	for _, user := range users {
		runner.CastBallotAs(user, testutil.MainChannelID, true /* inFavor */)
		runner.SendVoteStatusMessage(testutil.MainChannelID)
	}

	runner.ExpireVote(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)
}

func TestVote_Fail(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Initialize users.
	users := []*discordgo.User{
		testutil.NewUser("user0", 0 /* id */, false /* bot */),
		testutil.NewUser("user1", 1 /* id */, false /* bot */),
		testutil.NewUser("user2", 2 /* id */, false /* bot */),
		testutil.NewUser("user3", 3 /* id */, false /* bot */),
		testutil.NewUser("user4", 4 /* id */, false /* bot */),
	}
	for _, user := range users {
		runner.AddUser(user)
	}

	// Start the vote.
	runner.SendVoteMessageAs(users[0], testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Cast votes
	for _, user := range users {
		runner.CastBallotAs(user, testutil.MainChannelID, false /* inFavor */)
		runner.SendVoteStatusMessage(testutil.MainChannelID)
	}

	runner.ExpireVote(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)
}

func TestVote_Tie(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Initialize users.
	users := []*discordgo.User{
		testutil.NewUser("user0", 0 /* id */, false /* bot */),
		testutil.NewUser("user1", 1 /* id */, false /* bot */),
		testutil.NewUser("user2", 2 /* id */, false /* bot */),
		testutil.NewUser("user3", 3 /* id */, false /* bot */),
		testutil.NewUser("user4", 4 /* id */, false /* bot */),
		testutil.NewUser("user5", 5 /* id */, false /* bot */),
		testutil.NewUser("user6", 6 /* id */, false /* bot */),
		testutil.NewUser("user7", 7 /* id */, false /* bot */),
		testutil.NewUser("user8", 8 /* id */, false /* bot */),
		testutil.NewUser("user9", 9 /* id */, false /* bot */),
	}
	for _, user := range users {
		runner.AddUser(user)
	}

	// Start the vote.
	runner.SendVoteMessageAs(users[0], testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Cast votes
	for _, user := range users {
		runner.CastBallotAs(user, testutil.MainChannelID, false /* inFavor */)
		runner.SendVoteStatusMessage(testutil.MainChannelID)
	}

	runner.ExpireVote(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)
}

func TestVote_TwoVotes(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Initialize users.
	users := []*discordgo.User{
		testutil.NewUser("user0", 0 /* id */, false /* bot */),
		testutil.NewUser("user1", 1 /* id */, false /* bot */),
		testutil.NewUser("user2", 2 /* id */, false /* bot */),
		testutil.NewUser("user3", 3 /* id */, false /* bot */),
		testutil.NewUser("user4", 4 /* id */, false /* bot */),
	}
	for _, user := range users {
		runner.AddUser(user)
	}

	// Start the vote.
	runner.SendVoteMessageAs(users[0], testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Cast votes
	for _, user := range users {
		runner.CastBallotAs(user, testutil.MainChannelID, true /* inFavor */)
		runner.SendVoteStatusMessage(testutil.MainChannelID)
	}

	runner.ExpireVote(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Start the vote again.
	runner.SendVoteMessageAs(users[0], testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Cast votes again.
	for _, user := range users {
		runner.CastBallotAs(user, testutil.MainChannelID, false /* inFavor */)
		runner.SendVoteStatusMessage(testutil.MainChannelID)
	}

	runner.ExpireVote(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)
}

func TestVote_TwoChannels(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Initialize users.
	users := []*discordgo.User{
		testutil.NewUser("user0", 0 /* id */, false /* bot */),
		testutil.NewUser("user1", 1 /* id */, false /* bot */),
		testutil.NewUser("user2", 2 /* id */, false /* bot */),
		testutil.NewUser("user3", 3 /* id */, false /* bot */),
		testutil.NewUser("user4", 4 /* id */, false /* bot */),
	}
	for _, user := range users {
		runner.AddUser(user)
	}

	// Start the votes.
	runner.SendVoteMessageAs(users[0], testutil.MainChannelID)

	// Stagger the vote starts, since in practice votes will never start exactly
	// at the same moment.
	runner.UTCClock.Advance(vote.VoteDuration / 2)
	runner.UTCTimer.ElapseTime(vote.VoteDuration / 2)

	runner.SendVoteMessageAs(users[0], testutil.SecondChannelID)

	runner.SendVoteStatusMessage(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.SecondChannelID)

	// Cast votes
	for _, user := range users {
		runner.CastBallotAs(user, testutil.MainChannelID, true /* inFavor */)
		runner.CastBallotAs(user, testutil.SecondChannelID, false /* inFavor */)
		runner.SendVoteStatusMessage(testutil.MainChannelID)
		runner.SendVoteStatusMessage(testutil.SecondChannelID)
	}

	// Expires both votes.
	runner.ExpireVote(testutil.MainChannelID)
	runner.ExpireVote(testutil.SecondChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.SecondChannelID)
}

func TestVote_CannotVoteTwice(t *testing.T) {
	runner := testutil.NewRunner(t)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Calls vote with no args, and then actually starts a vote.
	author := testutil.NewUser("author", 0 /* id */, false /* bot */)
	runner.AddUser(author)
	runner.SendVoteMessageAs(author, testutil.MainChannelID)
	runner.CastBallotAs(author, testutil.MainChannelID, true /* inFavor */)
	runner.CastDuplicateBallotAs(author, testutil.MainChannelID, true /* inFavor */)
	runner.CastDuplicateBallotAs(author, testutil.MainChannelID, false /* inFavor */)
}
