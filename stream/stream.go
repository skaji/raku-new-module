package stream

import (
	"bytes"
	"context"
	"log"
	"net/mail"
	"time"

	"github.com/skaji/perl6-cpan-new/distribution"
	"github.com/skaji/perl6-cpan-new/nntp"
)

func fixPerl6Distribution(ctx context.Context, d *distribution.Distribution) {
	fetcher := distribution.NewPerl6Fetcher()
	max := 10
	for i := 1; i <= max; i++ {
		name, err := fetcher.FetchName(ctx, d)
		if err == nil {
			if d.MainModule == name {
				log.Printf("guessed MainModule %s matches name in META file", d.MainModule)
			} else {
				log.Printf("guessed MainModule %s does NOT match name (%s) in META file, use the name in META file", d.MainModule, name)
				d.MainModule = name
			}
			break
		}
		log.Print(err)
		if _, ok := err.(*distribution.RetryableError); !ok {
			break
		}
		if i != max {
			log.Println("Sleep 30sec...")
			time.Sleep(30 * time.Second)
		}
	}
}

// NewPerl6 is
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
			msg, err := mail.ReadMessage(bytes.NewReader(article.Article))
			if err != nil {
				log.Print(err)
				continue
			}
			subject := msg.Header.Get("Subject")
			dist, err := distribution.New(subject)
			if err != nil {
				log.Print(err)
				continue
			}
			log.Print(article.ID, " ", dist.AsJSON())
			if dist.IsPerl6 {
				go func() {
					fixPerl6Distribution(ctx, dist)
					ch <- dist
				}()
			}
		}
	}()
	return ch
}
