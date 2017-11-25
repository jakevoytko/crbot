package main

import (
	"fmt"
	"testing"

	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/feature/help"
	"github.com/jakevoytko/crbot/feature/learn"
	"github.com/jakevoytko/crbot/feature/list"
	"github.com/jakevoytko/crbot/testutil"
)

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
	runner.SendMessage(testutil.DirectMessageID, "?unlearn call", fmt.Sprintf(app.MsgPublicOnly, "?unlearn"))
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
}
