package network

import (
	"smartlog/client"
	"smartlog/uri"
)

func New(ur *uri.URI) (*client.Client, error) {
	c := &client.Client{
		URI: ur,
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}
	return c, nil
}
