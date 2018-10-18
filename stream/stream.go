package stream

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/mail"
	"time"

	"github.com/skaji/perl6-cpan-new/distribution"
	"github.com/skaji/perl6-cpan-new/nntp"
)

func fixPerl6Distribution(ctx context.Context, d *distribution.Distribution) error {
	fetcher := distribution.NewPerl6Fetcher()
	max := 20
	for i := 1; i <= max; i++ {
		name, err := fetcher.FetchName(ctx, d.MetaURL())
		if err == nil {
			if d.MainModule == name {
				log.Printf("%s guessed MainModule %s matches name in META file", d.ID, d.MainModule)
			} else {
				log.Printf("%s guessed MainModule %s does NOT match name (%s) in META file, use the name in META file", d.ID, d.MainModule, name)
				d.MainModule = name
			}
			return nil
		}
		if _, ok := err.(*distribution.RetryableError); !ok {
			return err
		}

		log.Println(d.ID, err)
		if i != max {
			log.Println(d.ID, "Sleep 30sec...")
			time.Sleep(30 * time.Second)
		}
	}
	return errors.New("too many retry, give up")
}

func NewPerl6(ctx context.Context, host string, port int, tick int) <-chan *distribution.Distribution {
	ch := make(chan *distribution.Distribution)
	go func() {
		nntpClient, err := nntp.NewClient(host, port, "perl.cpan.uploads", tick)
		if err != nil {
			log.Fatal(err)
		}
		defer nntpClient.Close()

		nntpChannel := nntpClient.Tail(ctx)

		for article := range nntpChannel {
			id := article.ID
			msg, err := mail.ReadMessage(bytes.NewReader(article.Article))
			if err != nil {
				log.Println(id, err)
				continue
			}
			subject := msg.Header.Get("Subject")
			dist, err := distribution.New(id, subject)
			if err != nil {
				log.Println(id, err)
				continue
			}

			log.Println(id, dist.AsJSON())
			if !dist.IsPerl6 {
				continue
			}

			go func(id string) {
				err := fixPerl6Distribution(ctx, dist)
				if err == nil {
					ch <- dist
				} else {
					log.Println(id, err)
				}
			}(id)
		}
	}()
	return ch
}
