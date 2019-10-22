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

type SlackLogger struct {
	url    string
	ch     chan string
	Logger Logger
}

func NewSlack(url string) Logger {
	l := &SlackLogger{
		url: url,
		ch:  make(chan string, 1000),
		Logger: &CoreLogger{
			Level:  4,
			Logger: base.New(os.Stderr, "", base.LstdFlags|base.Llongfile),
		},
	}
	go l.poster()
	return l
}

func (l *SlackLogger) Fatal(v ...interface{}) {
	l.Logger.Fatal(v...)
}

func (l *SlackLogger) Printf(format string, v ...interface{}) {
	l.Post(fmt.Sprintf(format, v...))
	l.Logger.Printf(format, v...)
}

func (l *SlackLogger) Println(v ...interface{}) {
	l.Post(fmt.Sprintln(v...))
	l.Logger.Println(v...)
}

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
