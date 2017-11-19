package vote

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/testutil"
)

func TestVote(t *testing.T) {
	runner := testutil.NewRunner(t)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Calls vote with no args, and then actually starts a vote.
	author := testutil.NewUser("author", 0 /* id */, false /* bot */)
	runner.AddUser(author)
	runner.SendMessageAs(author, testutil.MainChannelID, "?vote", MsgHelpVote)
	runner.SendVoteMessageAs(author, testutil.MainChannelID)
	runner.SendVoteStatusMessage(testutil.MainChannelID)

	// Assert that a second vote can't be started.
	runner.SendMessageAs(author, testutil.MainChannelID, "?vote another vote", MsgActiveVote)

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
	runner.UTCClock.Advance(VoteDuration / 2)
	runner.UTCTimer.ElapseTime(VoteDuration / 2)

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
