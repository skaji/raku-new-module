package twitter

import (
	orig "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Twitter interface {
	Tweet(string) error
}

type Client struct {
	*orig.Client
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string) *Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return &Client{orig.NewClient(httpClient)}
}

func (c *Client) Tweet(str string) error {
	_, _, err := c.Statuses.Update(str, nil)
	return err
}

type NoopClient struct{}

func NewNoop() *NoopClient {
	return &NoopClient{}
}

func (c *NoopClient) Tweet(string) error {
	return nil
}
