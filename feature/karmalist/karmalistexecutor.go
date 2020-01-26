package karmalist

import (
	"bytes"
	"sort"
	"strconv"
	"math"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

const (
	// MsgGistAddress is a user-visible string announcing the url of the hastebin
	MsgGistAddress = "The list of karma is here"
	// MsgListKarma is a user-visible header for the list of Karma'd things
	MsgListKarma = "Karma targets listed by intensity:"
)

// Executor uploads all command keys to hastebin and returns the url to the user
type Executor struct {
	featureRegistry *feature.Registry
	karmaMap      stringmap.StringMap
	gist            api.Gist
}

// NewExecutor works as advertised
func NewExecutor(featureRegistry *feature.Registry, karmaMap stringmap.StringMap, gist api.Gist) *Executor {
	return &Executor{
		featureRegistry: featureRegistry,
		karmaMap:        karmaMap,
		gist:            gist,
	}
}

// GetType returns the type of this feature.
func (e *Executor) GetType() int {
	return model.CommandTypeKarmaList
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *Executor) PublicOnly() bool {
	return false
}

// Execute uploads the command list to github and pings the gist link in chat.
func (e *Executor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	all, err := e.karmaMap.GetAll()
	if err != nil {
		log.Fatal("Error reading all the karma", err)
	}

	type sortableKarma struct {
		displayKarma   string
		absKarma int
	}
	var karmaStore []sortableKarma

	for k, v := range all {
			displayKarma := k + ": " + v
			floatKarma, _ := strconv.ParseFloat(v, 32)
			absKarma := int(math.Abs(floatKarma))
			karmaStore = append(karmaStore, sortableKarma{displayKarma, absKarma})
	}

	sort.Slice(karmaStore, func(i, j int) bool {
			return karmaStore[i].absKarma > karmaStore[j].absKarma
	})

	var buffer bytes.Buffer
	buffer.WriteString(MsgListKarma)
	buffer.WriteString("\n")
	for _, kv := range karmaStore {
		buffer.WriteString(kv.displayKarma)
		buffer.WriteString("\n")
	}

	url, err := e.gist.Upload(buffer.String())
	if err != nil {
		s.ChannelMessageSend(channel.Format(), err.Error())
		return
	}
	s.ChannelMessageSend(channel.Format(), MsgGistAddress+": "+url)

}
