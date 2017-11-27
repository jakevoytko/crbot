package model

import "github.com/bwmarrin/discordgo"

///////////////////////////////////////////////////////////////////////////////
// Constants
///////////////////////////////////////////////////////////////////////////////

const (
	Type_Custom = iota
	Type_Help
	Type_Karma
	Type_Learn
	Type_List
	Type_None
	Type_RickList
	Type_RickListInfo
	Type_Unlearn
	Type_Unrecognized
	Type_Vote
	Type_VoteBallot
	Type_VoteConclude
	Type_VoteStatus

	Name_Help           = "?help"
	Name_KarmaIncrement = "?++"
	Name_KarmaDecrement = "?--"
	Name_Learn          = "?learn"
	Name_List           = "?list"
	Name_RickListInfo   = "?ricklist"
	Name_Unlearn        = "?unlearn"
	Name_Vote           = "?vote"
	Name_VoteAgainstF2  = "?f2"
	Name_VoteAgainstNo  = "?no"
	Name_VoteInFavorF1  = "?f1"
	Name_VoteInFavorYes = "?yes"
	Name_VoteStatus     = "?votestatus"
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

type VoteData struct {
	Message string
}

type BallotData struct {
	InFavor bool
}

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
