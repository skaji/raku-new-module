package nntp

import (
	"bytes"
	"context"
	"net/mail"
	"testing"
	"time"

	"github.com/skaji/perl6-cpan-new/pkg/log"
)

func TestNNTP(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := Tail(ctx, "nntp.perl.org:119", "perl.cpan.uploads", 2*time.Second, -5)
	count := 0
	for article := range ch {
		message, err := mail.ReadMessage(bytes.NewReader(article.Article))
		if err != nil {
			continue
		}
		log.Println(message.Header.Get("Subject"))
		count++
	}
	if count < 5 {
		t.Fatal("fail")
	}
}
