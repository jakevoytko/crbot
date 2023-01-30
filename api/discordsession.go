package api

import "github.com/bwmarrin/discordgo"

// DiscordSession is an interface for interacting with Discord within a session
// message handler.
type DiscordSession interface {
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	Channel(channelID string, options ...discordgo.RequestOption) (*discordgo.Channel, error)
	User(userID string, options ...discordgo.RequestOption) (*discordgo.User, error)
}
