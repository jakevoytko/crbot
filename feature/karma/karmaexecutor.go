package karma

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

type KarmaExecutor struct {
	karmaMap stringmap.StringMap
}

// NewKarmaExecutor works as advertised.
func NewKarmaExecutor(karmaMap stringmap.StringMap) *KarmaExecutor {
	return &KarmaExecutor{karmaMap: karmaMap}
}

// GetType returns the type of this feature.
func (e *KarmaExecutor) GetType() int {
	return model.Type_Karma
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *KarmaExecutor) PublicOnly() bool {
	return true
}

const (
	MsgIncrementKarma = "%v has been upvoted. %v now has %d karma."
	MsgDecrementKarma = "%v has been downvoted. %v now has %d karma."
)

// Execute attempts to add karma to the total already in memory, or creates a
// new record if it was not found.
func (e *KarmaExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	if command.Karma == nil {
		log.Fatal("Incorrectly generated karma command", errors.New("wat"))
	}

	// Get the current value of karma (if it exists) and increment/decrement it
	currentKarma := 0
	has, err := e.karmaMap.Has(command.Karma.Target)
	if err != nil {
		log.Fatal("Couldn't check if target has karma", err)
	}
	if has {
		currentKarmaStr, err := e.karmaMap.Get(command.Karma.Target)
		if err != nil {
			log.Fatal("Couldn't get target's current karma", err)
		}
		currentKarma, err = strconv.Atoi(currentKarmaStr)
		if err != nil {
			log.Fatal("Invalid karma value", err)
		}
	}

	var newKarma int
	if command.Karma.Increment {
		newKarma = currentKarma + 1
	} else {
		newKarma = currentKarma - 1
	}

	err = e.karmaMap.Set(command.Karma.Target, strconv.Itoa(newKarma))
	if err != nil {
		log.Fatal("Error storing new karma value", err)
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
