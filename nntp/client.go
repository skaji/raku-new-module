package nntp

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/dustin/go-nntp/client"
)

// Article is
type Article struct {
	Article []byte
	ID      string
}

// Client is
type Client struct {
	Group      string
	Host       string
	Port       int
	Tick       time.Duration
	Offset     int64 // where to read
	Reconnect  bool
	backend    *nntpclient.Client
	currentID  int64
	previousID int64
}

// NewClient is
func NewClient(host string, port int, group string, tick int) (*Client, error) {
	backend, err := connect(host, port, group)
	if err != nil {
		return nil, err
	}
	if tick == 0 {
		tick = 30
	}
	return &Client{
		backend:    backend,
		Host:       host,
		Port:       port,
		Group:      group,
		Tick:       time.Duration(tick) * time.Second,
		Offset:     0,
		Reconnect:  true,
		previousID: -1,
		currentID:  -1,
	}, nil
}

func connect(host string, port int, group string) (*nntpclient.Client, error) {
	backend, err := nntpclient.New("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	// first time check
	if _, err := backend.Group(group); err != nil {
		return nil, err
	}
	return backend, nil
}

func (c *Client) reconnect() error {
	c.Close()
	backend, err := connect(c.Host, c.Port, c.Group)
	if err != nil {
		return err
	}
	c.backend = backend
	return nil
}

// Close is
func (c *Client) Close() error {
	c.previousID = -1
	c.currentID = -1
	return c.backend.Close()
}

// Tail is
func (c *Client) Tail(ctx context.Context) <-chan *Article {
	ch := make(chan *Article)
	go func() {
		ticker := time.NewTicker(c.Tick)
		defer ticker.Stop()

		reconnect := false
		for {
			select {
			case <-ticker.C:
				if reconnect && c.Reconnect {
					log.Println("try to reconnect...")
					if err := c.reconnect(); err != nil {
						log.Println("NG reconnect, ", err)
						continue
					}
					log.Println("OK reconnect")
					reconnect = false
					// pass through
				}
				group, err := c.backend.Group(c.Group)
				if err != nil {
					log.Println(err)
					reconnect = true
					continue
				}
				c.currentID = group.High
				if c.previousID == -1 {
					if c.currentID+c.Offset >= 0 {
						c.previousID = c.currentID + c.Offset
					} else {
						c.previousID = 0
					}
				}

				for i := c.previousID + 1; i <= c.currentID; i++ {
					var err error
					var r io.Reader
					var article []byte

					ID := strconv.FormatInt(i, 10)
					_, _, r, err = c.backend.Article(ID)
					if err == nil {
						article, err = ioutil.ReadAll(r)
						if err == nil {
							ch <- &Article{Article: article, ID: ID}
							continue
						}
					}
					log.Println(err)
					reconnect = true
				}
				c.previousID = c.currentID
			case <-ctx.Done():
				close(ch)
				return
			}
		}
	}()
	return ch
}
