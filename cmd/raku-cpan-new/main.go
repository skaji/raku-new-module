package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/skaji/raku-cpan-new/pkg/config"
	"github.com/skaji/raku-cpan-new/pkg/log"
	"github.com/skaji/raku-cpan-new/pkg/stream"
	"github.com/skaji/raku-cpan-new/pkg/twitter"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: raku-cpan-new config.json")
		os.Exit(1)
	}
	c, err := config.NewFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if c.SlackURL != "" {
		log.Set(log.NewSlack(c.SlackURL))
	}
	defer log.Close()

	log.Println("start")
	run(c)
	log.Println("finish")
}

func run(c *config.Config) {
	tw := twitter.NewNoop()
	if c.ConsumerKey != "" {
		log.Println("will tweet with ConsumerKey", c.ConsumerKey)
		tw = twitter.New(c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessSecret)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		s := <-sig
		log.Printf("catch %v\n", s)
		cancel()
	}()

	strm := stream.NewRaku(ctx, c.Addr, time.Duration(c.Tick)*time.Second)
	for dist := range strm {
		summary := dist.Summary()
		log.Println(dist.ID, "tweet", strings.Replace(summary, "\n", " ", -1))
		if err := tw.Tweet(summary); err != nil {
			log.Println(dist.ID, err)
		}
	}
}