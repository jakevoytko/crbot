package moderation

import (
	"testing"

	"github.com/aetimmes/discordgo"
	"github.com/jakevoytko/crbot/feature/help"
	"github.com/jakevoytko/crbot/feature/moderation"
	"github.com/jakevoytko/crbot/testutil"
)

func TestRickList(t *testing.T) {
	runner := testutil.NewRunner(t)

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

	// A non-learn message gets intercepted.
	runner.SendMessageAs(rickListedUser, testutil.DirectMessageID, "?help help-arg", moderation.MsgRickList)

	// A learn can still go through.
	runner.SendLearnMessageAs(rickListedUser, testutil.DirectMessageID, "?learn rick list", testutil.NewLearnData("rick", "list"))
}
