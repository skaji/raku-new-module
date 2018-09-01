package main

import (
	"bufio"
	"container/ring"
	"errors"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	nntp "github.com/dustin/go-nntp"
	"github.com/dustin/go-nntp/server"
)

// Backend is
type Backend struct {
	Group    *nntp.Group
	Articles []*nntp.Article
}

// type Article struct {
// 	// The article's headers
// 	Header textproto.MIMEHeader
// 	// The article's body
// 	Body io.Reader
// 	// Number of bytes in the article body (used by OVER/XOVER)
// 	Bytes int
// 	// Number of lines in the article body (used by OVER/XOVER)
// 	Lines int
// }

func newArticle(subject string) *nntp.Article {
	header := textproto.MIMEHeader{}
	header.Set("Subject", subject)
	body := "this is body\n"
	article := nntp.Article{
		Header: header,
		Body:   strings.NewReader(body),
		Bytes:  len(body),
		Lines:  1,
	}
	return &article
}

func emit(b *Backend) {
	file, err := os.Open("./cpan.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	lines := []string{}
	s := bufio.NewScanner(file)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	r := ring.New(len(lines))
	for i := 0; i < len(lines); i++ {
		r.Value = lines[i]
		r = r.Next()
	}
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			v := r.Value.(string)
			r = r.Next()
			b.Articles = append(b.Articles, newArticle(v))
			b.Group.Count++
			b.Group.High++
			if b.Group.Low == -1 {
				b.Group.Low++
			}
		}
	}
}

// NewBackend is
func NewBackend() *Backend {
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
	go emit(backend)
	return backend
}

// ListGroups is
func (b *Backend) ListGroups(max int) ([]*nntp.Group, error) {
	return nil, nil
}

// GetGroup is
func (b *Backend) GetGroup(name string) (*nntp.Group, error) {
	if name != b.Group.Name {
		return nil, nntpserver.ErrNoSuchGroup
	}
	return b.Group, nil
}

// GetArticle is
func (b *Backend) GetArticle(group *nntp.Group, id string) (*nntp.Article, error) {
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

// GetArticles is
func (b *Backend) GetArticles(group *nntp.Group, from, to int64) ([]nntpserver.NumberedArticle, error) {
	return nil, nil
}

// Authorized is
func (b *Backend) Authorized() bool {
	return true
}

// Authenticate is
func (b *Backend) Authenticate(user, pass string) (nntpserver.Backend, error) {
	return nil, nil
}

// AllowPost is
func (b *Backend) AllowPost() bool {
	return true
}

// Post is
func (b *Backend) Post(article *nntp.Article) error {
	return nil
}

func main() {
	a, err := net.ResolveTCPAddr("tcp", ":1119")
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	s := nntpserver.NewServer(NewBackend())
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.Process(c)
	}
}
