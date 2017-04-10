package util

import (
	"errors"

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
	Channels  map[string]*discordgo.Channel
	currentID int
	author    *discordgo.User
}

// NewInMemoryDiscordSession works as advertised.
func NewInMemoryDiscordSession() *InMemoryDiscordSession {
	channels := make(map[string]*discordgo.Channel)
	return &InMemoryDiscordSession{
		Messages:  []*Message{},
		Channels:  channels,
		currentID: 0,
		author: &discordgo.User{
			"bot_id",
			"bot@email.com",
			"crbot",
			"bot avatar",
			"bot discriminator",
			"bot token",
			true,  /* verified */
			false, /* multifactor enabled */
			true,  /* Bot */
		},
	}
}

// ChannelMessageSend records a new message delivery.
func (s *InMemoryDiscordSession) ChannelMessageSend(channel, message string) (*discordgo.Message, error) {
	s.Messages = append(s.Messages, NewMessage(channel, message))
	return &discordgo.Message{
		ID:              "id",
		ChannelID:       channel,
		Content:         message,
		Timestamp:       "timestamp",
		EditedTimestamp: "edited timestamp",
		MentionRoles:    []string{},
		Tts:             false,
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
func (s *InMemoryDiscordSession) Channel(channelID string) (*discordgo.Channel, error) {
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
