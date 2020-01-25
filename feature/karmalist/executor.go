package karmalist

import (
	"bytes"
	"sort"
	"fmt"
	//"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

const (
	// MsgGistAddress is a user-visible string announcing the url of the gist
	MsgGistAddress = "The list of karma is here"
	// MsgListBuiltins is a user-visible header for the list of builtins
	MsgListKarma = "Ahoy! Thar be karma below!"
	// MsgListCustom is a user-visible header for the list of learned commands
	MsgListCustom = "List of learned commands:"
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
		karmaMap:      karmaMap,
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
		log.Fatal("Error reading all karma", err)
	}
	// https://stackoverflow.com/questions/42629541/go-lang-sort-a-2d-array
	custom := make([]string, 0, len(all))
	for name, kcnt := range all {
		combo := name + ": " + kcnt
		custom = append(custom, combo)
	}
	sort.Strings(custom)
	var buffer bytes.Buffer
	// buffer.WriteString(strings.Join(all, ","))
	buffer.WriteString(MsgListKarma)
	buffer.WriteString("\n")
	for _, name := range custom {
		buffer.WriteString(name)
		fmt.Printf("%s\n", name )
		// if strings.Contains(all[name], "$1") {
		buffer.WriteString("\n")
		// }
		buffer.WriteString("\n")
	}
/*
	url, err := e.gist.Upload(buffer.String())
	if err != nil {
		s.ChannelMessageSend(channel.Format(), err.Error())
		return
	}
	s.ChannelMessageSend(channel.Format(), MsgGistAddress+": "+url)
	*/
}
