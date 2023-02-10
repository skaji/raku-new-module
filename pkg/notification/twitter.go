package notification

import (
	"context"

	orig "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Twitter struct {
	*orig.Client
}

func NewTwitter(consumerKey, consumerSecret, accessToken, accessSecret string) *Twitter {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return &Twitter{orig.NewClient(httpClient)}
}

func (t *Twitter) Notify(ctx context.Context, message string) error {
	_, _, err := t.Statuses.Update(message, nil)
	return err
}

var _ Notifier = (*Twitter)(nil)
