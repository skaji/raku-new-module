package log

import (
	"bytes"
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
	client *http.Client
	Logger Logger

	stop chan<- struct{}
	done <-chan struct{}
}

func NewSlack(url string) Logger {
	l := &SlackLogger{
		url:    url,
		ch:     make(chan string, 1000),
		client: &http.Client{Timeout: 3 * time.Second},
		Logger: &CoreLogger{
			Level:  4,
			Logger: base.New(os.Stderr, "", base.LstdFlags|base.Llongfile),
		},
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		l.poster(stop)
		close(done)
	}()
	l.stop = stop
	l.done = done
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

func (l *SlackLogger) Close() {
	close(l.stop)
	defer func() {
		l.Logger.Close()
		<-l.done
	}()
	for {
		select {
		case text := <-l.ch:
			if err := l.post(text); err != nil {
				l.Logger.Println(err)
			}
		default:
			return
		}
	}
}

func (l *SlackLogger) Post(text string) {
	select {
	case l.ch <- text:
	default:
		l.Logger.Println("slack channel is full, skip", text)
	}
}

func (l *SlackLogger) poster(stop <-chan struct{}) {
	for {
		select {
		case text := <-l.ch:
			if err := l.post(text); err != nil {
				l.Logger.Println(err)
			}
		case <-stop:
			return
		}

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

	res, err := l.client.Do(req)
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
