package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	redis "gopkg.in/redis.v5"
)

///////////////////////////////////////////////////////////////////////////////
// CRBot is a call-and-response bot. It is taught to learn a call and
// response. When it sees the call, it replays the response. Look at the ?help
// documentation for a full list of commands.
//
// Licensed under MIT license, at project root.
///////////////////////////////////////////////////////////////////////////////

func main() {
	var filename = flag.String("filename", "secret.json", "Filename of configuration json")
	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		fatal("Could not ping Redis", err)
	}

	secret, e := ParseSecret(*filename)
	if e != nil {
		fatal("Secret parsing failed", e)
	}

	discord, err := discordgo.New("Bot " + secret.BotToken)
	if err != nil {
		fatal("Error initializing Discord client library", e)
	}

	discord.AddHandler(getHandleMessage(redisClient))

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
	Type_None = iota
	Type_Unrecognized
	Type_Help
	Type_Learn
	Type_Custom
	Type_List

	Name_Help  = "?help"
	Name_Learn = "?learn"
	Name_List  = "?list"

	Redis_Hash = "crbot-custom-commands"
)

// TypeToString maps builtin types to their string names.
var TypeToString map[int]string = map[int]string{
	Type_Help:  Name_Help,
	Type_Learn: Name_Learn,
	Type_List:  Name_List,
}

// StringToType maps builtin names to their types.
var StringToType map[string]int = map[string]int{
	Name_Help:  Type_Help,
	Name_Learn: Type_Learn,
	Name_List:  Type_List,
}

// Tries to get a value from s and ?s
func getUserStringToType(s string) int {
	if t, ok := StringToType[s]; ok {
		return t
	} else if t, ok := StringToType["?"+s]; ok {
		return t
	}
	return Type_None
}

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

// getHandleMessage returns a function to parse and handle incoming messages.
func getHandleMessage(redisClient *redis.Client) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Never reply to a bot.
		if m.Author.Bot {
			return
		}

		command, err := parseCommand(redisClient, m.Content)
		if err != nil {
			info("Error processing command", err)
		}

		switch command.Type {
		case Type_Help:
			sendHelp(s, m.ChannelID, command)
		case Type_Learn:
			sendLearn(redisClient, s, m.ChannelID, command)
		case Type_Custom:
			sendCustom(redisClient, s, m.ChannelID, command)
		case Type_List:
			sendList(redisClient, s, m.ChannelID, command)
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

type CustomData struct {
	Call string
	Args string
}

type Command struct {
	Type   int
	Help   *HelpData
	Learn  *LearnData
	Custom *CustomData
}

// Parses the raw text string from the user. Returns an executable command.
func parseCommand(redisClient *redis.Client, content string) (*Command, error) {
	if !strings.HasPrefix(content, "?") {
		return &Command{
			Type: Type_None,
		}, nil
	}
	splitContent := strings.Split(content, " ")

	// Parse builtins.
	switch splitContent[0] {
	case Name_Help:
		return parseHelp(splitContent)
	case Name_Learn:
		return parseLearn(redisClient, splitContent)
	case Name_List:
		return &Command{Type: Type_List}, nil
	}

	// See if it's a custom command.
	if redisClient.HExists(Redis_Hash, splitContent[0][1:]).Val() {
		return parseCustom(redisClient, splitContent)
	}

	// No such command!
	return &Command{
		Type: Type_Unrecognized,
	}, nil
}

func parseHelp(splitContent []string) (*Command, error) {
	if splitContent[0] != Name_Help {
		fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	userType := Type_Unrecognized
	if len(splitContent) > 1 {
		userType = getUserStringToType(splitContent[1])
	}
	return &Command{
		Type: Type_Help,
		Help: &HelpData{
			Type: userType,
		},
	}, nil
}

func parseLearn(redisClient *redis.Client, splitContent []string) (*Command, error) {
	if splitContent[0] != Name_Learn {
		fatal("parseLearn called with non-learn command", errors.New("wat"))
	}

	callRegexp := regexp.MustCompile("(?s)^[[:alnum:]].*$")
	responseRegexp := regexp.MustCompile("(?s)^[^/?!].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 3 || !callRegexp.MatchString(splitContent[1]) || !responseRegexp.MatchString(splitContent[2]) {
		return &Command{
			Type: Type_Help,
			Help: &HelpData{
				Type: Type_Learn,
			},
		}, nil
	}

	// Don't overwrite old or builtin commands.
	if redisClient.HExists(Redis_Hash, splitContent[1]).Val() || getUserStringToType(splitContent[1]) != Type_None {
		return &Command{
			Type: Type_Learn,
			Learn: &LearnData{
				CallOpen: false,
				Call:     splitContent[1],
			},
		}, nil
	}

	// Everything is good.
	response := strings.Join(splitContent[2:], " ")
	return &Command{
		Type: Type_Learn,
		Learn: &LearnData{
			CallOpen: true,
			Call:     splitContent[1],
			Response: response,
		},
	}, nil
}

func parseCustom(redisClient *redis.Client, splitContent []string) (*Command, error) {
	if !redisClient.HExists(Redis_Hash, splitContent[0][1:]).Val() {
		fatal("parseCustom called with non-custom command", errors.New("wat"))
	}
	return &Command{
		Type: Type_Custom,
		Custom: &CustomData{
			Call: splitContent[0][1:],
			Args: strings.Join(splitContent[1:], " "),
		},
	}, nil
}

///////////////////////////////////////////////////////////////////////////////
// User-visible messages
///////////////////////////////////////////////////////////////////////////////

const (
	MsgCustomNeedsArgs   = "This command takes args. Please type `?command <more text>` instead of `?command`"
	MsgDefaultHelp       = "Type `?help` for this message, `?list` to list all commands, or `?help <command>` to get help for a particular command."
	MsgGistAddress       = "The list of commands is here: "
	MsgGistPostFail      = "Unable to connect to Gist service. Give it a few minutes and try again"
	MsgGistResponseFail  = "Failure reading response from Gist service"
	MsgGistSerializeFail = "Unable to serialize Gist"
	MsgGistStatusCode    = "Failed to upload Gist :("
	MsgGistUrlFail       = "Failed getting url from Gist service"
	MsgHelpHelp          = "You're probably right. I probably didn't think of this case."
	MsgHelpLearn         = "Type `?learn <call> <the response the bot should read>`. When you type `?call`, the bot will reply with the response.\n\nThe first character of the call must be alphanumeric, and the first character of the response must not begin with /, ?, or !\n\nUse $1 in the response to substitute all arguments"
	MsgHelpList          = "Type `?list` to get the URL of a Gist with all builtin and learned commands"
	MsgLearnFail         = "I already know ?%s"
	MsgLearnSuccess      = "Learned about %s"
	MsgListBuiltins      = "List of builtins:"
	MsgListCustom        = "List of learned commands:"
)

///////////////////////////////////////////////////////////////////////////////
// Response
///////////////////////////////////////////////////////////////////////////////

func sendHelp(s *discordgo.Session, channel string, command *Command) {
	if command.Help == nil {
		fatal("Incorrectly generated help command", errors.New("wat"))
	}
	switch command.Help.Type {
	default:
		if _, err := s.ChannelMessageSend(channel, MsgDefaultHelp); err != nil {
			info("Failed to send default help message", err)
		}
	case Type_Help:
		s.ChannelMessageSend(channel, MsgHelpHelp)
	case Type_Learn:
		s.ChannelMessageSend(channel, MsgHelpLearn)
	case Type_List:
		s.ChannelMessageSend(channel, MsgHelpList)
	}
}

func sendLearn(redisClient *redis.Client, s *discordgo.Session, channel string, command *Command) {
	if command.Learn == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}
	if !command.Learn.CallOpen {
		s.ChannelMessageSend(channel, fmt.Sprintf(MsgLearnFail, command.Learn.Call))
		return
	}

	// Teach the command.
	if redisClient.HExists(Redis_Hash, command.Learn.Call).Val() {
		fatal("Collision when adding a call for "+command.Learn.Call, errors.New("wat"))
	}
	redisClient.HSet(Redis_Hash, command.Learn.Call, command.Learn.Response)

	// Send ack.
	s.ChannelMessageSend(channel, fmt.Sprintf(MsgLearnSuccess, command.Learn.Call))
}

func sendCustom(redisClient *redis.Client, s *discordgo.Session, channel string, command *Command) {
	if command.Custom == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}

	if !redisClient.HExists(Redis_Hash, command.Custom.Call).Val() {
		fatal("Accidentally found a mismatched call/response pair", errors.New("Call response mismatch"))
	}

	response := redisClient.HGet(Redis_Hash, command.Custom.Call).Val()

	if strings.Contains(response, "$1") {
		if command.Custom.Args == "" {
			s.ChannelMessageSend(channel, MsgCustomNeedsArgs)
			return
		}
		response = strings.Replace(response, "$1", command.Custom.Args, 1)
	}
	s.ChannelMessageSend(channel, response)
}

func sendList(redisClient *redis.Client, s *discordgo.Session, channel string, command *Command) {
	builtins := []string{}
	for name := range StringToType {
		builtins = append(builtins, name)
	}

	custom := []string{}
	for name := range redisClient.HGetAll(Redis_Hash).Val() {
		custom = append(custom, name)
	}

	sort.Strings(builtins)
	sort.Strings(custom)

	var buffer bytes.Buffer
	buffer.WriteString(MsgListBuiltins)
	buffer.WriteString("\n")
	for _, name := range builtins {
		buffer.WriteString(" - ")
		buffer.WriteString(name)
		buffer.WriteString("\n")
	}

	buffer.WriteString("\n")

	buffer.WriteString(MsgListCustom)
	buffer.WriteString("\n")
	for _, name := range custom {
		buffer.WriteString(" - ?")
		buffer.WriteString(name)
		buffer.WriteString("\n")
	}

	url, err := uploadCommandList(buffer.String())
	if err != nil {
		s.ChannelMessageSend(channel, err.Error())
		return
	}
	s.ChannelMessageSend(channel, MsgGistAddress+": "+url)
}

///////////////////////////////////////////////////////////////////////////////
// Gist handling
///////////////////////////////////////////////////////////////////////////////
type Gist struct {
	Description string           `json:"description"`
	Public      bool             `json:"public"`
	Files       map[string]*File `json:"files"`
}

// A file represents the contents of a Gist.
type File struct {
	Content string `json:"content"`
}

// simpleGist returns a Gist object with just the given contents.
func simpleGist(contents string) *Gist {
	return &Gist{
		Public:      false,
		Description: "CRBot command list",
		Files: map[string]*File{
			"commands": &File{
				Content: contents,
			},
		},
	}
}

func uploadCommandList(contents string) (string, error) {
	gist := simpleGist(contents)
	serializedGist, err := json.Marshal(gist)
	if err != nil {
		info("Error marshalling gist", err)
		return "", errors.New(MsgGistSerializeFail)
	}
	response, err := http.Post(
		"https://api.github.com/gists", "application/json", bytes.NewBuffer(serializedGist))
	if err != nil {
		info("Error POSTing gist", err)
		return "", errors.New(MsgGistPostFail)
	} else if response.StatusCode != 201 {
		info("Bad status code", errors.New("Code: "+strconv.Itoa(response.StatusCode)))
		return "", errors.New(MsgGistStatusCode)
	}

	responseMap := map[string]interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		info("Error reading gist response", err)
		return "", errors.New(MsgGistResponseFail)
	}

	if finalUrl, ok := responseMap["html_url"]; ok {
		return finalUrl.(string), nil
	}
	return "", errors.New(MsgGistUrlFail)
}
