package app

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jakevoytko/crbot/model"
)

///////////////////////////////////////////////////////////////////////////////
// Configuration handling
///////////////////////////////////////////////////////////////////////////////

// Secret holds the serialized bot token.
type Config struct {
	BotToken string            `json:"bot_token"`
	RickList []model.Snowflake `json:"ricklist"`
}

// ParseConfig reads the config from the given filename.
func ParseConfig(filename string) (*Config, error) {
	f, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	var config Config
	e = json.Unmarshal(f, &config)
	return &config, e
}
