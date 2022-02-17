package none

import (
	"github.com/KarelKubat/smartlog/client"
	"github.com/KarelKubat/smartlog/uri"
)

func New(ur *uri.URI) (*client.Client, error) {
	return &client.Client{
		URI: ur,
	}, nil
}
