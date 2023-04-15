package notification

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Mastodon struct {
	URL         string
	AccessToken string
}

func NewMastodon(u string, accessToken string) *Mastodon {
	return &Mastodon{URL: u, AccessToken: accessToken}
}

func (m *Mastodon) Notify(ctx context.Context, message string) error {
	data := url.Values{}
	data.Set("status", message)
	return m.post(ctx, m.URL+"/api/v1/statuses", data)
}

func (m *Mastodon) post(ctx context.Context, u string, data url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+m.AccessToken)
	res, err := http.DefaultClient.Do(req)
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

var _ Notifier = (*Mastodon)(nil)
