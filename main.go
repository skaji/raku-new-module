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
						log.Print(r.Err)
						return
					}
					log.Println(" producer")
					log.Printf("  %v\n", r.Distribution)
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
