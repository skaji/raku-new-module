package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dghubble/oauth1"
)

type Twitter struct {
	client *http.Client
}

func NewTwitter(consumerKey, consumerSecret, accessToken, accessSecret string) *Twitter {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	return &Twitter{client: config.Client(oauth1.NoContext, token)}
}

func (t *Twitter) Notify(ctx context.Context, message string) error {
	const tweetURL = "https://api.twitter.com/2/tweets"

	b, err := json.Marshal(map[string]string{"text": message})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tweetURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := t.client.Do(req)
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, res.Body)
	res.Body.Close()
	if res.StatusCode/100 != 2 {
		return errors.New(res.Status)
	}
	return nil
}

var _ Notifier = (*Twitter)(nil)
