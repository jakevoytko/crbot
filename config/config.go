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

// NewConfig builds a new config and sets default values for config params that have them.
func NewConfig() Config {
	return Config{
		RedisHost:     "localhost",
		RedisPort:     6379,
		RedisPassword: "",
		RedisDatabase: 0,
	}
}

// ParseConfig reads the config from the given filename.
func ParseConfig(filename string) (*Config, error) {
	f, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	config := NewConfig()
	e = json.Unmarshal(f, &config)
	return &config, e
}
