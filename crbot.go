package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
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
	config, err := app.ParseConfig(*filename)
	if err != nil {
		log.Fatal("Config parsing failed", err)
	}

	// Initialize Redis.
	commandMap, err := NewRedisStringMap(Redis_Hash)
	if err != nil {
		log.Fatal("Unable to initialize Redis", err)
	}
	gist := api.NewRemoteGist()

	// Set up Discord API.
	discord, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal("Error initializing Discord client library", err)
	}

	featureRegistry := InitializeRegistry(commandMap, gist, config)

	// Open communications with Discord.
	handler := getHandleMessage(commandMap, featureRegistry)

	// Wrapper is needed so the discordgo registry recognizes the input types.
	wrappedHandler := func(s *discordgo.Session, c *discordgo.MessageCreate) {
		handler(s, c)
	}
	discord.AddHandler(wrappedHandler)
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord session", err)
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
