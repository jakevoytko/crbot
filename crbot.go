package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

///////////////////////////////////////////////////////////////////////////////
// CRBot is a call-and-response bot. It is taught by users to learn a call and
// response. When it sees the call, it replays the response. Look at the ?help
// documentation for a full list of commands.
//
// Licensed under MIT license, at project root.
///////////////////////////////////////////////////////////////////////////////

func main() {
	var filename = flag.String("filename", "secret.json", "Filename of configuration json")
	flag.Parse()

	// Parse config.
	secret, err := ParseSecret(*filename)
	if err != nil {
		fatal("Secret parsing failed", err)
	}

	// TODO(jake): Refactor Features to provide multiple parsers and executors,
	// and add this to the Learn feature.
	commandMap, err := NewRedisStringMap(Redis_Hash)
	if err != nil {
		fatal("Unable to initialize Redis", err)
	}
	gist := NewRemoteGist()

	featureRegistry := InitializeRegistry(commandMap, gist)

	// Set up Discord API.
	discord, err := discordgo.New("Bot " + secret.BotToken)
	if err != nil {
		fatal("Error initializing Discord client library", err)
	}

	// Open communications with Discord.
	discord.AddHandler(getHandleMessage(commandMap, featureRegistry))
	if err := discord.Open(); err != nil {
		fatal("Error opening Discord session", err)
	}

	fmt.Println("CRBot running.")

	<-make(chan interface{})
}

///////////////////////////////////////////////////////////////////////////////
// Constants
///////////////////////////////////////////////////////////////////////////////

const (
	Redis_Hash = "crbot-custom-commands"
)

///////////////////////////////////////////////////////////////////////////////
// Configuration handling
///////////////////////////////////////////////////////////////////////////////

// Secret holds the serialized bot token.
type Secret struct {
	BotToken string `json:"bot_token"`
}

// ParseSecret reads the config from the given filename.
func ParseSecret(filename string) (*Secret, error) {
	f, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	var config Secret
	e = json.Unmarshal(f, &config)
	return &config, e
}
