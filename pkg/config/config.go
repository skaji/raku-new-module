package config

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ConsumerKey    string `json:"consumer_key" envconfig:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret" envconfig:"consumer_secret"`
	AccessToken    string `json:"access_token" envconfig:"access_token"`
	AccessSecret   string `json:"access_secret" envconfig:"access_secret"`
	Addr           string `json:"addr" envconfig:"addr"`
	Tick           int    `json:"tick" envconfig:"tick"`
	SlackURL       string `json:"slack_url" envconfig:"slack_url"`
	DiscordURL     string `json:"discord_url" envconfig:"discord_url"`
}

func NewFromFile(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func NewFromEnv() (*Config, error) {
	var c Config
	if err := envconfig.Process("app", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
