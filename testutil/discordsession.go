package testutil

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Message encapsulates a channel/message pair.
type Message struct {
	Channel string
	Message string
}

// NewMessage works as advertised.
func NewMessage(channel, message string) *Message {
	return &Message{
		Channel: channel,
		Message: message,
	}
}

// InMemoryDiscordSession is a fake for the discord session.
type InMemoryDiscordSession struct {
	Messages  []*Message
	Users     map[string]*discordgo.User
	Channels  map[string]*discordgo.Channel
	currentID int
	author    *discordgo.User
}

// NewInMemoryDiscordSession works as advertised.
func NewInMemoryDiscordSession() *InMemoryDiscordSession {
	channels := make(map[string]*discordgo.Channel)

	users := make(map[string]*discordgo.User)

	author := &discordgo.User{
		ID:            "bot_id",
		Email:         "bot@email.com",
		Username:      "crbot",
		Avatar:        "bot avatar",
		Discriminator: "bot discriminator",
		Token:         "bot token",
		Verified:      true,  /* verified */
		MFAEnabled:    false, /* multifactor enabled */
		Bot:           true,  /* Bot */
	}

	rickListedUser := &discordgo.User{
		ID:            "2",
		Email:         "bot@email.com",
		Username:      "crbot",
		Avatar:        "bot avatar",
		Discriminator: "bot discriminator",
		Token:         "bot token",
		Verified:      true,  /* verified */
		MFAEnabled:    false, /* multifactor enabled */
		Bot:           false, /* Bot */
	}

	users["bot_id"] = author
	users["2"] = rickListedUser

	return &InMemoryDiscordSession{
		Messages:  []*Message{},
		Channels:  channels,
		Users:     users,
		currentID: 0,
		author:    author,
	}
}

// ChannelMessageSend records a new message delivery.
func (s *InMemoryDiscordSession) ChannelMessageSend(channel string, message string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	s.Messages = append(s.Messages, NewMessage(channel, message))
	editedTimestamp := time.Now()
	return &discordgo.Message{
		ID:              "id",
		ChannelID:       channel,
		Content:         message,
		Timestamp:       time.Now().Add(-time.Hour),
		EditedTimestamp: &editedTimestamp,
		MentionRoles:    []string{},
		TTS:             false,
		MentionEveryone: false,
		Author:          s.author,
		Attachments:     []*discordgo.MessageAttachment{},
		Embeds:          []*discordgo.MessageEmbed{},
		Mentions:        []*discordgo.User{},
		Reactions:       []*discordgo.MessageReactions{},
	}, nil
}

// Channel returns the Channel struct of the given channel ID. Can be used to
// determine attributes such as the channel name, topic, etc.
func (s *InMemoryDiscordSession) Channel(channelID string, options ...discordgo.RequestOption) (*discordgo.Channel, error) {
	if channel := s.Channels[channelID]; channel != nil {
		return channel, nil
	}
	return nil, errors.New("Attempted to get missing channel " + channelID)
}

// SetChannel adds a channel to the map of channels that InMemoryDiscordSession
// holds. Can be used to test features that only work under certain channel
// conditions.
func (s *InMemoryDiscordSession) SetChannel(channel *discordgo.Channel) {
	s.Channels[channel.ID] = channel
}

// User returns the user struct of the given user ID.
func (s *InMemoryDiscordSession) User(userID string, options ...discordgo.RequestOption) (*discordgo.User, error) {
	if user := s.Users[userID]; user != nil {
		return user, nil
	}
	return nil, errors.New("Attempted to get missing user " + userID)
}
