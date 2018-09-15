package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jakevoytko/crbot/model"
)

///////////////////////////////////////////////////////////////////////////////
// Configuration handling
///////////////////////////////////////////////////////////////////////////////

// Config represents the JSON format of the config file.
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
