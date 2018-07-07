package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"time"

	"github.com/dustin/go-nntp/client"
)

// NNTP is
type NNTP struct {
	Group      string
	Server     string
	Tick       time.Duration
	Offset     int64
	currentID  int64
	previousID int64
}

// NewNNTP is
func NewNNTP(server string, group string) *NNTP {
	return &NNTP{
		Group:      group,
		Server:     server,
		Tick:       30 * time.Second,
		Offset:     0,
		currentID:  -1,
		previousID: -1,
	}
}

// Result is
type Result struct {
	Distribution *Distribution
	Err          error
}

// Tail is
func (n *NNTP) Tail(ctx context.Context) <-chan *Result {
	ch := make(chan *Result)
	go func() {
		client, err := nntpclient.New("tcp", fmt.Sprintf("%s:%d", n.Server, 119))
		if err != nil {
			ch <- &Result{Err: err}
			close(ch)
			return
		}
		group, err := client.Group(n.Group)
		if err != nil {
			ch <- &Result{Err: err}
		}

		readBody := func(client *nntpclient.Client, ID int64) (*Distribution, error) {
			_, _, r, err := client.Article(strconv.FormatInt(ID, 10))
			if err != nil {
				return nil, err
			}
			m, err := mail.ReadMessage(r)
			if err != nil {
				return nil, err
			}
			distribution, err := NewDistribution(m.Header.Get("Subject"))
			if err != nil {
				return nil, err
			}
			return distribution, nil
		}

		ticker := time.NewTicker(n.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				group, err = client.Group(n.Group)
				if err != nil {
					ch <- &Result{Err: err}
					continue
				}
				n.currentID = group.High
				if n.previousID == -1 {
					n.previousID = n.currentID + n.Offset
				}

				seen := 0
				for i := n.previousID + 1; i <= n.currentID; i++ {
					distribution, err := readBody(client, i)
					ch <- &Result{Distribution: distribution, Err: err}
					seen++
					if seen > 20 {
						log.Printf("  seen more than 20 articles, break")
						break
					}
				}
				n.previousID = n.currentID
			case <-ctx.Done():
				log.Println(" ctx.Done()")
				client.Close()
				close(ch)
				return
			}
		}
	}()
	return ch
}
