package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-nntp/client"
)

var distributionRegexp = regexp.MustCompile(`\$CPAN/(authors/id/./../([^/]+)/(Perl6/)?(.+)\.(?:tar\.gz|tar\.bz2|zip|tgz))`)

func parseBody(body string) (*Distribution, error) {
	res := distributionRegexp.FindAllStringSubmatch(body, -1)
	if len(res) == 0 {
		return nil, errors.New("failed to parse")
	}

	r := res[0]
	d := Distribution{
		CPANID:    r[2],
		Distvname: r[4],
		IsPerl6:   false,
		Pathname:  r[1],
	}
	if r[3] == "Perl6/" {
		d.IsPerl6 = true
	}
	parts := strings.Split(d.Distvname, "-")
	d.Distname = strings.Join(parts[:len(parts)-1], "-")
	d.MainModule = strings.Join(parts[:len(parts)-1], "::")
	return &d, nil
}

type NNTP struct {
	CurrentID  int64
	Group      string
	PreviousID int64
	Server     string
	Tick       time.Duration
}

func NewNNTP(server string, group string) *NNTP {
	return &NNTP{
		CurrentID:  -1,
		Group:      group,
		PreviousID: -1,
		Server:     server,
		Tick:       30 * time.Second,
	}
}

type Result struct {
	Distribution *Distribution
	Err          error
}

func (n *NNTP) Tail(ctx context.Context) <-chan *Result {
	ch := make(chan *Result)
	go func() {
		client, err := nntpclient.New("tcp", fmt.Sprintf("%s:%d", n.Server, 119))
		if err != nil {
			ch <- &Result{Err: err}
			return
		}
		group, err := client.Group(n.Group)
		if err != nil {
			ch <- &Result{Err: err}
		}

		readBody := func(client *nntpclient.Client, ID int64) (*Distribution, error) {
			_, _, r, err := client.Body(strconv.FormatInt(ID, 10))
			if err != nil {
				return nil, err
			}
			body, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
			distribution, err := parseBody(string(body))
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
				n.CurrentID = group.High
				if n.PreviousID == -1 {
					n.PreviousID = n.CurrentID
				}

				seen := 0
				for i := n.PreviousID + 1; i <= n.CurrentID; i++ {
					distribution, err := readBody(client, i)
					ch <- &Result{Distribution: distribution, Err: err}
					seen++
					if seen > 5 {
						log.Printf("  seen more than 5 articles, break")
						break
					}
				}
				n.PreviousID = n.CurrentID
			case <-ctx.Done():
				log.Println(" ctx.Done()")
				client.Close()
				return
			}
		}
	}()
	return ch
}
