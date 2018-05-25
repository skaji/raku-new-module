package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	log.Println("start")
	twitter, err := NewTwitter("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	for {
		done := false
		func() {
			nntp := NewNNTP("nntp.perl.org", "perl.cpan.uploads")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			produce := nntp.Tail(ctx)
			for {
				select {
				case r := <-produce:
					if r.Err != nil {
						log.Println(r.Err)
						if _, ok := r.Err.(*DistributionError); ok {
							continue
						} else {
							return
						}
					}
					log.Println(r.Distribution.AsJSON())
					if !r.Distribution.IsPerl6 {
						continue
					}
					_, _, err := twitter.Statuses.Update(r.Distribution.Summary(), nil)
					if err != nil {
						log.Println(err)
						continue
					}
				case s := <-sig:
					log.Printf("catch %v\n", s)
					done = true
					return
				}

			}
		}()
		if done {
			break
		}
		log.Println("Retry after 60sec...")
		time.Sleep(60 * time.Second)
	}
	log.Println("finish")
}
