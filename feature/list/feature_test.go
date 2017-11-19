package list

import (
	"testing"

	"github.com/jakevoytko/crbot/testutil"
)

func TestGarbage_NoResponse(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Commands that should never return a response.
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "?")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "!")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, ".")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "!help")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, "help")
	runner.SendMessageWithoutResponse(testutil.MainChannelID, ".help")

	// Assert that the list functionality is OK after a bunch of garbage commands.
	runner.SendListMessage(testutil.MainChannelID)
}
