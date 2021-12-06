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
	BotToken      string            `json:"bot_token"`
	RickList      []model.Snowflake `json:"ricklist"`
	RedisHost     string            `json:"redis_host"`
	RedisPort     int               `json:"redis_port"`
	RedisUsername string            `json:"redis_username"`
	RedisPassword string            `json:"redis_password"`
	RedisDatabase int               `json:"redis_database"`
}

// SetDefaultConfig sets default values for config params that have them.
func SetDefaultConfig(config *Config) {
	config.RedisHost = "localhost"
	config.RedisPort = 6379
	config.RedisPassword = ""
	config.RedisDatabase = 0
}

// ParseConfig reads the config from the given filename.
func ParseConfig(filename string) (*Config, error) {
	f, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	var config Config
	SetDefaultConfig(&config)
	e = json.Unmarshal(f, &config)
	return &config, e
}
