package config

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// also known as API Key
	TwitterConsumerKey string `json:"twitter_consumer_key" envconfig:"twitter_consumer_key"`
	// also known as API Key Secret
	TwitterConsumerSecret string `json:"twitter_consumer_secret" envconfig:"twitter_consumer_secret"`

	TwitterAccessToken  string `json:"twitter_access_token" envconfig:"twitter_access_token"`
	TwitterAccessSecret string `json:"twitter_access_secret" envconfig:"twitter_access_secret"`
	MastodonURL         string `json:"mastodon_url" envconfig:"mastodon_url"`
	MastodonAccessToken string `json:"mastodon_access_token" envconfig:"mastodon_access_token"`
	RecentURL           string `json:"recent_url" envconfig:"recent_url"`
	Tick                int    `json:"tick" envconfig:"tick"`
	SlackURL            string `json:"slack_url" envconfig:"slack_url"`
	DiscordURL          string `json:"discord_url" envconfig:"discord_url"`
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
