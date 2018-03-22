package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	IFTTTKey string `json:"ifttt_key"`
}

func LoadConfig(file string) (*Config, error) {
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
