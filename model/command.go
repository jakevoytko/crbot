package model

import "github.com/bwmarrin/discordgo"

///////////////////////////////////////////////////////////////////////////////
// Constants
///////////////////////////////////////////////////////////////////////////////

// Consts use throughout the application
const (
	CommandTypeCustom = iota
	CommandTypeFactSphere
	CommandTypeHelp
	CommandTypeKarma
	CommandTypeKarmaList
	CommandTypeLearn
	CommandTypeList
	CommandTypeNone
	CommandTypeRickList
	CommandTypeRickListInfo
	CommandTypeUnlearn
	CommandTypeUnrecognized
	CommandTypeVote
	CommandTypeVoteBallot
	CommandTypeVoteConclude
	CommandTypeVoteStatus

	CommandNameFactSphere     = "?factsphere"
	CommandNameHelp           = "?help"
	CommandNameKarmaIncrement = "?++"
	CommandNameKarmaDecrement = "?--"
	CommandNameKarmaList      = "?karmalist"
	CommandNameLearn          = "?learn"
	CommandNameList           = "?list"
	CommandNameRickListInfo   = "?ricklist"
	CommandNameUnlearn        = "?unlearn"
	CommandNameVote           = "?vote"
	CommandNameVoteAgainstF2  = "?f2"
	CommandNameVoteAgainstNo  = "?no"
	CommandNameVoteInFavorF1  = "?f1"
	CommandNameVoteInFavorYes = "?yes"
	CommandNameVoteStatus     = "?votestatus"
)

///////////////////////////////////////////////////////////////////////////////
// User message parsing
///////////////////////////////////////////////////////////////////////////////

// HelpData holds data for Help commands.
type HelpData struct {
	Command string
}

// KarmaData holds the target and whether karma is to be incremented or
// decremented
type KarmaData struct {
	Increment bool
	Target    string
}

// LearnData is the learn-specific data
type LearnData struct {
	CallOpen bool
	Call     string
	Response string
}

// UnlearnData is the unlearn-specific data
type UnlearnData struct {
	CallOpen bool
	Call     string
}

// CustomData is the custom ?learn-specific data
type CustomData struct {
	Call string
	Args string
}

// VoteData contains the information about the proposed vote
type VoteData struct {
	Message string
}

// BallotData represents whether the user is for or against the vote
type BallotData struct {
	InFavor bool
}

// Command is the generic command interface
// TODO(jake): Make this an interface that has only getType(), cast in features.
type Command struct {
	// Metadata
	Author       *discordgo.User
	ChannelID    Snowflake
	Type         int
	OriginalName string

	// Message data
	Ballot  *BallotData
	Custom  *CustomData
	Help    *HelpData
	Karma   *KarmaData
	Learn   *LearnData
	Unlearn *UnlearnData
	Vote    *VoteData
}
