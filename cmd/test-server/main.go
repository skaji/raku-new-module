package main

import (
	"context"
	"log"
	"net"

	nntpserver "github.com/dustin/go-nntp/server"
)

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

	s := nntpserver.NewServer(NewBackend(context.Background()))
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.Process(c)
	}
}
