package model

import "github.com/bwmarrin/discordgo"

///////////////////////////////////////////////////////////////////////////////
// Constants
///////////////////////////////////////////////////////////////////////////////

const (
	Type_Custom = iota
	Type_Help
	Type_Learn
	Type_List
	Type_None
	Type_RickList
	Type_RickListInfo
	Type_Unlearn
	Type_Unrecognized
	Type_VoteStatus
	Type_Vote

	Name_Help         = "?help"
	Name_Learn        = "?learn"
	Name_List         = "?list"
	Name_Unlearn      = "?unlearn"
	Name_RickListInfo = "?ricklist"
	Name_Vote         = "?vote"
	Name_VoteStatus   = "?votestatus"
)

///////////////////////////////////////////////////////////////////////////////
// User message parsing
///////////////////////////////////////////////////////////////////////////////

// HelpData holds data for Help commands.
type HelpData struct {
	Command string
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

// TODO(jake): Make this an interface that has only getType(), cast in features.
type Command struct {
	Custom  *CustomData
	Help    *HelpData
	Learn   *LearnData
	Type    int
	Unlearn *UnlearnData
	Author  *discordgo.User
	Vote    *VoteData
}
