package nntp

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	nntpclient "github.com/skaji/go-nntp/client"
	"github.com/skaji/raku-cpan-new/pkg/log"
)

var (
	aLongTimeAgo = time.Unix(1, 0)
	noTimeout    = time.Time{}
)

type Article struct {
	Article []byte
	ID      string
}

func Tail(globalCtx context.Context, addr string, group string, tick time.Duration, offset int64) <-chan *Article {
	ch := make(chan *Article)

	go func() {
		defer close(ch)

		var (
			client     *nntpclient.Client
			previousID int64 = -1
			currentID  int64
		)

		oneTick := func(ctx context.Context) error {
			if client == nil {
				var err error
				client, err = nntpclient.New(ctx, "tcp", addr)
				if err != nil {
					return fmt.Errorf("nntp connect: %w", err)
				}
				log.Println("connect OK")
			}

			cancel := client.SetContext(ctx)
			defer cancel()
			g, err := client.Group(group)
			if err != nil {
				return fmt.Errorf("nntp read group: %w", err)
			}
			currentID = g.High
			log.Debugf("%d %s high", g.High, group)
			if previousID == -1 {
				previousID = currentID + offset
			}
			for i := previousID + 1; i <= currentID; i++ {
				ID := strconv.FormatInt(i, 10)
				_, _, r, err := client.Article(ID)
				if err != nil {
					return fmt.Errorf("nntp read article: %w", err)
				}
				article, err := ioutil.ReadAll(r)
				if err != nil {
					return fmt.Errorf("nntp read article body: %w", err)
				}
				ch <- &Article{Article: article, ID: ID}
			}
			previousID = currentID
			return nil
		}
		reset := func() {
			if client != nil {
				log.Println("close connection")
				client.Close()
				client = nil
			}
			previousID = -1
		}
		defer reset()

		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
			err := oneTick(ctx)
			cancel()
			if err != nil {
				log.Println(err)
				reset()
			}
			select {
			case <-ticker.C:
			case <-globalCtx.Done():
				return
			}
		}
	}()
	return ch

}
