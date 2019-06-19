package twitter

import (
	orig "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client is
type Client struct {
	*orig.Client
}

// New is
func New(consumerKey, consumerSecret, accessToken, accessSecret string) *Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return &Client{orig.NewClient(httpClient)}
}

// Tweet is
func (c *Client) Tweet(str string) error {
	_, _, err := c.Statuses.Update(str, nil)
	return err
}
