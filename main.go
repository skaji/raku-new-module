package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("start")

	c, err := LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	ifttt := NewIFTTT("perl6_cpan_new", c.IFTTTKey)
	nntp := NewNNTP("nntp.perl.org", "perl.cpan.uploads")

	for {
		done := false
		func() {
			sig := make(chan os.Signal)
			signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			produce := nntp.Tail(ctx)
			for {
				select {
				case r := <-produce:
					if r.Err != nil {
						log.Println(r.Err)
						continue
					}
					log.Printf("  %s\n", r.Distribution.AsJSON())
					// if !r.Distribution.IsPerl6 {
					// 	continue
					// }
					err := ifttt.Post(r.Distribution.Summary())
					if err != nil {
						log.Println(err)
						continue
					}
				case s := <-sig:
					log.Printf(" catch %v\n", s)
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
