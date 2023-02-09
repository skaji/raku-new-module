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
	"github.com/skaji/raku-new-module/pkg/notification"
	"github.com/skaji/raku-new-module/pkg/stream"
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
	var notifiers notification.Notifiers
	if c.TwitterConsumerKey != "" {
		log.Print("will notify to twitter")
		notifiers = append(notifiers, notification.NewTwitter(
			c.TwitterConsumerKey, c.TwitterConsumerSecret, c.TwitterAccessToken, c.TwitterAccessSecret,
		))
	}
	if c.MastodonAccessToken != "" {
		log.Print("will notify to mastodon")
		notifiers = append(notifiers, notification.NewMastodon(
			c.MastodonURL, c.MastodonAccessToken,
		))
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

		log.Print(dist.ID, "notify", dist.URL)
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		if err := notifiers.Notify(ctx, dist.Summary()); err != nil {
			log.Print(dist.ID, err)
		}
		cancel()
	}
	return nil
}
