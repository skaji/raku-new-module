package main

import (
	"container/ring"
	"context"
	"errors"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"time"

	nntp "github.com/dustin/go-nntp"
	nntpserver "github.com/dustin/go-nntp/server"
)

func newArticle(subject string) *nntp.Article {
	header := textproto.MIMEHeader{}
	header.Set("Subject", subject)
	body := "this is body\n"
	return &nntp.Article{
		Header: header,
		Body:   strings.NewReader(body),
		Bytes:  len(body),
		Lines:  1,
	}
}

type Backend struct {
	Group    *nntp.Group
	Articles []*nntp.Article
	mux      sync.Mutex
}

func NewBackend(ctx context.Context) *Backend {
	backend := &Backend{
		Group: &nntp.Group{
			Name:        "perl.cpan.uploads",
			Description: "hoge",
			Count:       0,
			High:        -1,
			Low:         -1,
			Posting:     nntp.PostingNotPermitted,
		},
	}
	go backend.post(ctx)
	return backend
}

func (b *Backend) post(ctx context.Context) {
	lines := strings.Split(cpanText, "\n")
	r := ring.New(len(lines))
	for i := 0; i < len(lines); i++ {
		r.Value = lines[i]
		r = r.Next()
	}

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			v := r.Value.(string)
			r = r.Next()
			b.mux.Lock()
			b.Articles = append(b.Articles, newArticle(v))
			b.Group.Count++
			b.Group.High++
			if b.Group.Low == -1 {
				b.Group.Low++
			}
			b.mux.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (b *Backend) ListGroups(max int) ([]*nntp.Group, error) {
	return nil, nil
}

func (b *Backend) GetGroup(name string) (*nntp.Group, error) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if name != b.Group.Name {
		return nil, nntpserver.ErrNoSuchGroup
	}
	return b.Group, nil
}

func (b *Backend) GetArticle(group *nntp.Group, id string) (*nntp.Article, error) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if group.Name != b.Group.Name {
		return nil, nntpserver.ErrNoSuchGroup
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if len(b.Articles) <= i {
		return nil, errors.New("invalid id " + id)
	}
	return b.Articles[i], nil
}

func (b *Backend) GetArticles(group *nntp.Group, from, to int64) ([]nntpserver.NumberedArticle, error) {
	return nil, nil
}

func (b *Backend) Authorized() bool {
	return true
}

func (b *Backend) Authenticate(user, pass string) (nntpserver.Backend, error) {
	return nil, nil
}

func (b *Backend) AllowPost() bool {
	return true
}

func (b *Backend) Post(article *nntp.Article) error {
	return nil
}
