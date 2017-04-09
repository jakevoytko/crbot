package util

import "github.com/bwmarrin/discordgo"

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
	currentID int
	author    *discordgo.User
}

// NewInMemoryDiscordSession works as advertised.
func NewInMemoryDiscordSession() *InMemoryDiscordSession {
	return &InMemoryDiscordSession{
		Messages:  []*Message{},
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
