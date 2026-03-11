package client

import (
	"github.com/lumi4x/news/src/services/auth/client"
)

type Client struct {
	c *client.Client
}

func NewClient(client *client.Client) *Client {
	return &Client{
		c: client,
	}
}
