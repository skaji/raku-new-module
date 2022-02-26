package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skaji/raku-new-module/pkg/config"
	"github.com/skaji/raku-new-module/pkg/log"
	"github.com/skaji/raku-new-module/pkg/stream"
	"github.com/skaji/raku-new-module/pkg/twitter"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: raku-new-module config.json/-config-from-env")
		os.Exit(1)
	}
	var (
		c   *config.Config
		err error
	)
	if os.Args[1] == "-config-from-env" {
		c, err = config.NewFromEnv()
	} else {
		c, err = config.NewFromFile(os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}
	if c.SlackURL != "" {
		log.Set(log.NewSlack(c.SlackURL))
	} else if c.DiscordURL != "" {
		log.Set(log.NewDiscord(c.DiscordURL))
	}
	defer log.Close()

	log.Print("start")
	if err := run(c); err != nil {
		log.Fatal(err)
	}
	log.Print("finish")
}

func run(c *config.Config) error {
	tw := twitter.NewNoop()
	if c.ConsumerKey != "" {
		log.Print("will tweet with ConsumerKey", c.ConsumerKey)
		tw = twitter.New(c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessSecret)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	strm, err := stream.NewRaku(ctx, c.RecentURL, time.Duration(c.Tick)*time.Second)
	if err != nil {
		return err
	}
	for dist := range strm {
		if !dist.IsZef() && !dist.IsCPAN() {
			log.Print(dist.ID, "skip")
			continue
		}
		log.Print(dist.ID, "tweet", dist.URL)
		if err := tw.Tweet(dist.Summary()); err != nil {
			log.Print(dist.ID, err)
		}
	}
	return nil
}
