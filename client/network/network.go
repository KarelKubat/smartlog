package network

import (
	"github.com/KarelKubat/smartlog/client"
	"github.com/KarelKubat/smartlog/uri"
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
