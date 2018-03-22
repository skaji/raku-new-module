package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Config struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	AccessSecret   string `json:"access_secret"`
}

func loadConfig(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	if err := json.Unmarshal(content, c); err != nil {
		return nil, err
	}
	return c, nil
}

func newTwitter(file string) (*twitter.Client, error) {
	c, err := loadConfig(file)
	if err != nil {
		return nil, err
	}
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient), nil
}
