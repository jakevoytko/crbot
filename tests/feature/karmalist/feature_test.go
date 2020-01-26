package karmalist

import (
	"testing"

	"github.com/jakevoytko/crbot/testutil"
)

func TestKarmaList_NoResponse(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Add Karma using builtin
	runner.SendMessageIgnoringResponse(testutil.MainChannelID, "?++ Testing")
	// Test Karmalist verifying builtin success by proxy
	runner.SendKarmaListMessage(testutil.MainChannelID)
}
