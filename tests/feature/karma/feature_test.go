package karma

import (
	"fmt"
	"testing"

	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/feature/karma"
	"github.com/jakevoytko/crbot/testutil"
)

func TestKarma(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Correct call format
	runner.SendMessage(testutil.MainChannelID, "?++ target", fmt.Sprintf(karma.MsgIncrementKarma, "target", "target", 1))
	runner.SendMessage(testutil.MainChannelID, "?-- target", fmt.Sprintf(karma.MsgDecrementKarma, "target", "target", 0))

	// Can't give karma in private message
	runner.SendMessage(testutil.DirectMessageID, "?++ target", fmt.Sprintf(app.MsgPublicOnly, "?++"))

	// Wrong call format
	runner.SendMessage(testutil.MainChannelID, "?++", karma.MsgHelpKarmaIncrement)
	runner.SendMessage(testutil.MainChannelID, "?--", karma.MsgHelpKarmaDecrement)
	runner.SendMessage(testutil.MainChannelID, "?++ ?call response", karma.MsgHelpKarmaIncrement)
	runner.SendMessage(testutil.MainChannelID, "?-- ?call response", karma.MsgHelpKarmaDecrement)
}

func TestKarmaIncrement(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Test that karma accumlates and that @ and # are stripped
	runner.SendMessage(testutil.MainChannelID, "?++ target", fmt.Sprintf(karma.MsgIncrementKarma, "target", "target", 1))
	runner.SendMessage(testutil.MainChannelID, "?++ @target", fmt.Sprintf(karma.MsgIncrementKarma, "target", "target", 2))
	runner.SendMessage(testutil.MainChannelID, "?++ @target#0491", fmt.Sprintf(karma.MsgIncrementKarma, "target", "target", 3))
}

func TestKarmaDecrement(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Test that karma accumlates and that @ and # are stripped
	runner.SendMessage(testutil.MainChannelID, "?-- target", fmt.Sprintf(karma.MsgDecrementKarma, "target", "target", -1))
	runner.SendMessage(testutil.MainChannelID, "?-- @target#999", fmt.Sprintf(karma.MsgDecrementKarma, "target", "target", -2))
	runner.SendMessage(testutil.MainChannelID, "?-- @target#1337", fmt.Sprintf(karma.MsgDecrementKarma, "target", "target", -3))

}
