package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	redis "gopkg.in/redis.v5"
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

	// Initialize redis.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		fatal("Could not ping Redis", err)
	}

	// Parse config.
	secret, e := ParseSecret(*filename)
	if e != nil {
		fatal("Secret parsing failed", e)
	}

	// Initializing builtin features.
	// TODO(jvoytko): investigate the circularity that emerged to see if there's
	// a better pattern here.
	featureRegistry := NewFeatureRegistry()
	featureRegistry.Register(NewHelpFeature(featureRegistry))
	featureRegistry.Register(NewLearnFeature(featureRegistry, redisClient))
	featureRegistry.Register(NewUnlearnFeature(featureRegistry, redisClient))
	featureRegistry.Register(NewListFeature(featureRegistry, redisClient))
	customFeature := NewCustomFeature(redisClient)
	featureRegistry.Register(customFeature)
	featureRegistry.FallbackFeature = customFeature

	// Set up Discord API.
	discord, err := discordgo.New("Bot " + secret.BotToken)
	if err != nil {
		fatal("Error initializing Discord client library", e)
	}

	// Open communications with Discord.
	discord.AddHandler(getHandleMessage(redisClient, featureRegistry))
	if err := discord.Open(); err != nil {
		fatal("Error opening Discord session", err)
	}

	fmt.Println("CRBot running.")

	<-make(chan interface{})
}

///////////////////////////////////////////////////////////////////////////////
// Utility methods
///////////////////////////////////////////////////////////////////////////////

// fatal handles a non-recoverable error.
func fatal(msg string, err error) {
	panic(msg + ": " + err.Error())
}

func info(msg string, err error) {
	fmt.Printf(msg+": %v\n", err.Error())
}

///////////////////////////////////////////////////////////////////////////////
// Constants
///////////////////////////////////////////////////////////////////////////////

const (
	Type_Custom = iota
	Type_Help
	Type_Learn
	Type_List
	Type_None
	Type_Unlearn
	Type_Unrecognized

	Name_Help    = "?help"
	Name_Learn   = "?learn"
	Name_List    = "?list"
	Name_Unlearn = "?unlearn"

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

///////////////////////////////////////////////////////////////////////////////
// Controller methods
///////////////////////////////////////////////////////////////////////////////

// getHandleMessage returns the main handler for incoming messages.
func getHandleMessage(redisClient *redis.Client, featureRegistry *FeatureRegistry) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Never reply to a bot.
		if m.Author.Bot {
			return
		}

		command, err := parseCommand(redisClient, featureRegistry, m.Content)
		if err != nil {
			info("Error parsing command", err)
			return
		}

		feature := featureRegistry.GetFeatureByType(command.Type)
		if feature != nil {
			feature.Execute(s, m.ChannelID, command)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// User message parsing
///////////////////////////////////////////////////////////////////////////////

// HelpData holds data for Help commands.
type HelpData struct {
	Type int
}

type LearnData struct {
	CallOpen bool
	Call     string
	Response string
}

type UnlearnData struct {
	CallOpen bool
	Call     string
}

type CustomData struct {
	Call string
	Args string
}

// TODO(jake): Make this an interface that has only getType(), cast in features.
type Command struct {
	Custom  *CustomData
	Help    *HelpData
	Learn   *LearnData
	Type    int
	Unlearn *UnlearnData
}

// Parses the raw text string from the user. Returns an executable command.
func parseCommand(redisClient *redis.Client, registry *FeatureRegistry, content string) (*Command, error) {
	if !strings.HasPrefix(content, "?") {
		return &Command{
			Type: Type_None,
		}, nil
	}
	splitContent := strings.Split(content, " ")

	// Parse builtins.
	if feature := registry.GetFeatureByName(splitContent[0]); feature != nil {
		return feature.Parse(splitContent)
	}

	// See if it's a custom command.
	if redisClient.HExists(Redis_Hash, splitContent[0][1:]).Val() {
		return registry.FallbackFeature.Parse(splitContent)
	}

	// No such command!
	return &Command{
		Type: Type_Unrecognized,
	}, nil
}
