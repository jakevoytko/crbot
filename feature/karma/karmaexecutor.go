package karma

import (
	"errors"
	"fmt"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// Executor increments or decrements karma and prints the results to the user.
type Executor struct {
	modelHelper *ModelHelper
}

// NewExecutor works as advertised.
func NewExecutor(modelHelper *ModelHelper) *Executor {
	return &Executor{modelHelper: modelHelper}
}

// GetType returns the type of this feature.
func (e *Executor) GetType() int {
	return model.CommandTypeKarma
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *Executor) PublicOnly() bool {
	return true
}

const (
	// MsgIncrementKarma prints the results of ?++ thing
	MsgIncrementKarma = "%v has been upvoted. %v now has %d karma."
	// MsgDecrementKarma prints the results of ?-- thing
	MsgDecrementKarma = "%v has been downvoted. %v now has %d karma."
)

// Execute attempts to add karma to the total already in memory, or creates a
// new record if it was not found.
func (e *Executor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	if command.Karma == nil {
		log.Fatal("Incorrectly generated karma command", errors.New("wat"))
	}

	var newKarma int
	var err error
	if command.Karma.Increment {
		newKarma, err = e.modelHelper.Increment(command.Karma.Target)
	} else {
		newKarma, err = e.modelHelper.Decrement(command.Karma.Target)
	}

	if err != nil {
		log.Fatal("Error writing karma storage", err)
	}

	// Send ack.
	karmaAckMessage := MsgDecrementKarma
	if command.Karma.Increment {
		karmaAckMessage = MsgIncrementKarma
	}
	_, err = s.ChannelMessageSend(channelID.Format(), fmt.Sprintf(karmaAckMessage, command.Karma.Target, command.Karma.Target, newKarma))
	if err != nil {
		log.Info("Error sending karma message", err)
	}
}
