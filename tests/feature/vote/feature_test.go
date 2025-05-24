package vote

import (
	"testing"

	"github.com/jakevoytko/crbot/testutil"
)

const DeprecatedMSG = "Deprecated. Please use Discord polls."

func TestVote(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Calls vote with no args, and then actually starts a vote.
	author := testutil.NewUser("author", 0 /* id */, false /* bot */)
	runner.AddUser(author)
	runner.SendMessageAs(author, testutil.MainChannelID, "?vote", DeprecatedMSG)
	runner.SendMessageAs(author, testutil.MainChannelID, "?votestatus", DeprecatedMSG)
}
