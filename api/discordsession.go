package api

import "github.com/bwmarrin/discordgo"

// DiscordSession is an interface for interacting with Discord within a session
// message handler.
type DiscordSession interface {
	ChannelMessageSend(channel, message string) (*discordgo.Message, error)
	Channel(channelID string) (*discordgo.Channel, error)
	User(userID string) (*discordgo.User, error)
}
