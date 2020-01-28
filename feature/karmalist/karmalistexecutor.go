package karmalist

import (
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

// Executor uploads sorted karma list to hastebin and returns the url to the user
type Executor struct {
	featureRegistry *feature.Registry
	modelHelper     *ModelHelper
	karmaMap        stringmap.StringMap
	gist            api.Gist
}

// NewExecutor works as advertised
func NewExecutor(featureRegistry *feature.Registry, karmaMap stringmap.StringMap, modelHelper *ModelHelper, gist api.Gist) *Executor {
	return &Executor{
		featureRegistry: featureRegistry,
		modelHelper:     modelHelper,
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

// Execute uploads the sorted karma list to the gist API and pings the gist link in chat.
func (e *Executor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	sortedKarma := e.modelHelper.GenerateList()
	if url, err := e.gist.Upload(sortedKarma); err != nil {
		s.ChannelMessageSend(channel.Format(), err.Error())
		log.Info("Gist API failed", err)
	} else {
		s.ChannelMessageSend(channel.Format(), MsgGistAddress+": "+url)
	}

}
