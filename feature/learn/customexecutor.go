package learn

import (
	"errors"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

type CustomExecutor struct {
	commandMap model.StringMap
}

func NewCustomExecutor(commandMap model.StringMap) *CustomExecutor {
	return &CustomExecutor{commandMap: commandMap}
}

// GetType returns the type of this feature.
func (e *CustomExecutor) GetType() int {
	return model.Type_Custom
}

// Execute returns the response if possible.
func (e *CustomExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	if command.Custom == nil {
		log.Fatal("Incorrectly generated learn command", errors.New("wat"))
	}

	has, err := e.commandMap.Has(command.Custom.Call)
	if err != nil {
		log.Fatal("Error testing custom feature", err)
	}
	if !has {
		log.Fatal("Accidentally found a mismatched call/response pair", errors.New("Call response mismatch"))
	}

	response, err := e.commandMap.Get(command.Custom.Call)
	if err != nil {
		log.Fatal("Error reading custom response", err)
	}

	if strings.Contains(response, "$1") {
		if command.Custom.Args == "" {
			response = MsgCustomNeedsArgs
		} else {
			response = strings.Replace(response, "$1", command.Custom.Args, 4)
		}
	}
	s.ChannelMessageSend(channel.Format(), response)
}
