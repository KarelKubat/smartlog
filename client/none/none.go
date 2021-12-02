package none

import (
	"smartlog/client"
	"smartlog/uri"
)

func New(ur *uri.URI) (*client.Client, error) {
	return &client.Client{
		URI: ur,
	}, nil
}
