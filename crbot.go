package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/config"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
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
	var localhost = flag.String("localhost", "localhost", "Configurable localhost hostname, to allow for Docker's weirdness on OSX")

	flag.Parse()

	// Parse config.
	config, err := config.ParseConfig(*filename)
	if err != nil {
		log.Fatal("Config parsing failed", err)
	}

	// Initialize redis.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     *localhost + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal("Unable to initialize Redis", err)
	}

	commandMap := model.NewRedisStringMap(redisClient, RedisCommandHash)
	voteMap := model.NewRedisStringMap(redisClient, RedisVoteHash)

	gist := api.NewRemoteGist()

	// Set up Discord API.
	discord, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal("Error initializing Discord client library", err)
	}

	clock := model.NewSystemUTCClock()
	timer := model.NewSystemUTCTimer()

	// A command channel large enough to process a few commands without needing to
	// block.
	commandChannel := make(chan *model.Command, 10)

	featureRegistry := app.InitializeRegistry(
		commandMap, voteMap, gist, config, clock, timer, commandChannel)

	// Run any initial load handlers up front.
	for _, fn := range featureRegistry.GetInitialLoadFns() {
		err := fn(discord)
		if err != nil {
			log.Info("Error running initial load function", err)
		}
	}

	go app.HandleCommands(featureRegistry, discord, commandChannel)

	// Open communications with Discord.
	handler := app.GetHandleMessage(commandMap, featureRegistry, commandChannel)

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

// NOTE: These cannot change without a migration, since they are mapped to storage.
const (
	RedisCommandHash = "crbot-custom-commands"
	RedisVoteHash    = "crbot-feature-vote"
)
