package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	base "log"
	"net/http"
	"os"
	"time"
)

// SlackLogger is
type SlackLogger struct {
	url string
	ch  chan string
	*base.Logger
}

// NewSlack is
func NewSlack(url string) Logger {
	l := &SlackLogger{
		url:    url,
		ch:     make(chan string, 1000),
		Logger: base.New(os.Stderr, "", base.LstdFlags|base.Llongfile),
	}
	go l.poster()
	return l
}

// Fatal is
func (l *SlackLogger) Fatal(v ...interface{}) {
	l.Logger.Fatal(v...)
}

// Printf is
func (l *SlackLogger) Printf(format string, v ...interface{}) {
	l.Post(fmt.Sprintf(format, v...))
	l.Logger.Printf(format, v...)
}

// Println is
func (l *SlackLogger) Println(v ...interface{}) {
	l.Post(fmt.Sprintln(v...))
	l.Logger.Println(v...)
}

// Post is
func (l *SlackLogger) Post(text string) {
	select {
	case l.ch <- text:
	default:
		l.Logger.Println("slack channel is full, skip", text)
	}
}

func (l *SlackLogger) poster() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		text := <-l.ch
		go func() {
			if err := l.post(text); err != nil {
				l.Logger.Println(err)
			}
		}()
	}
}

func (l *SlackLogger) post(text string) error {
	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, l.url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	io.Copy(ioutil.Discard, res.Body)
	if res.StatusCode == 200 {
		return nil
	}
	return errors.New(res.Status)
}
