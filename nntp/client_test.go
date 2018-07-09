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
	fail := 0
	for article := range ch {
		if article.Error != nil {
			fail++
			continue
		}
		message, err := mail.ReadMessage(bytes.NewReader(article.Article))
		if err != nil {
			fail++
			continue
		}
		fmt.Println(message.Header.Get("Subject"))
	}
	if fail > 1 {
		t.Fatal("fail")
	}
}
