package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	AccessSecret   string `json:"access_secret"`
	Addr           string `json:"addr"`
	Tick           int    `json:"tick"`
	SlackURL       string `json:"slack_url"`
}

func NewFromFile(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
