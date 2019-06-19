package nntp

import (
	"bytes"
	"context"
	"fmt"
	"net/mail"
	"testing"
	"time"
)

func TestNNTP(t *testing.T) {
	nntp, err := NewClient("nntp.perl.org", 119, "perl.cpan.uploads", 5)
	if err != nil {
		panic(err)
	}
	nntp.Offset = -5
	nntp.Tick = time.Second
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	ch := nntp.Tail(ctx)
	count := 0
	for article := range ch {
		message, err := mail.ReadMessage(bytes.NewReader(article.Article))
		if err != nil {
			continue
		}
		fmt.Println(message.Header.Get("Subject"))
		count++
	}
	if count < 5 {
		t.Fatal("fail")
	}
}
