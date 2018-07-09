package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/skaji/perl6-cpan-new/stream"
	"github.com/skaji/perl6-cpan-new/twitter"
)

type config struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	AccessSecret   string `json:"access_secret"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Tick           int    `json:"tick"`
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
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("start")
	run()
	log.Println("finish")
}

func run() {
	c, err := loadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	var tw *twitter.Client
	if c.ConsumerKey != "" {
		tw = twitter.New(c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessSecret)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	stream := stream.NewPerl6(ctx, c.Host, c.Port, c.Tick)

	for {
		select {
		case dist := <-stream:
			summary := dist.Summary()
			log.Print("tweet ", strings.Replace(summary, "\n", " ", -1))
			if tw != nil {
				err := tw.Tweet(summary)
				if err != nil {
					log.Println(err)
				}
			}
		case s := <-sig:
			log.Printf("catch %v\n", s)
			cancel()
			return
		}
	}
}
