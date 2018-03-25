package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNNTP(t *testing.T) {
	nntp := NewNNTP("nntp.perl.org", "perl.cpan.uploads")
	nntp.Offset = -2
	nntp.Tick = time.Second
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	ch := nntp.Tail(ctx)
	fail := 0
	for r := range ch {
		fmt.Println(r.Distribution.AsJSON())
		if r.Err != nil {
			fail++
		}
	}
	if fail > 1 {
		t.Fatal("faile to parse")
	}
}
