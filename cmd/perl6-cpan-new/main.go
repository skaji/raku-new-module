package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/skaji/perl6-cpan-new/pkg/log"
	"github.com/skaji/perl6-cpan-new/pkg/stream"
	"github.com/skaji/perl6-cpan-new/pkg/twitter"
)

type config struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	AccessSecret   string `json:"access_secret"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Tick           int    `json:"tick"`
	SlackURL       string `json:"slack_url"`
}

func loadConfig(file string) (*config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c config
	if err := json.Unmarshal(content, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func main() {
	configfile := "./_test_config.json"
	if len(os.Args) > 1 {
		configfile = os.Args[1]
	}
	c, err := loadConfig(configfile)
	if err != nil {
		log.Fatal(err)
	}
	if c.SlackURL != "" {
		log.Set(log.NewSlack(c.SlackURL))
	}

	log.Println("start")
	run(c)
	log.Println("finish")
}

func run(c *config) {
	var tw *twitter.Client
	if c.ConsumerKey != "" {
		log.Println("will tweet with ConsumerKey", c.ConsumerKey)
		tw = twitter.New(c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessSecret)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	strm := stream.NewPerl6(ctx, c.Host, c.Port, c.Tick)

	for {
		select {
		case dist := <-strm:
			summary := dist.Summary()
			log.Println(dist.ID, "tweet", strings.Replace(summary, "\n", " ", -1))
			if tw != nil {
				err := tw.Tweet(summary)
				if err != nil {
					log.Println(dist.ID, err)
				}
			}
		case s := <-sig:
			log.Printf("catch %v\n", s)
			return
		}
	}
}
