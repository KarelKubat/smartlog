package any

import (
	"smartlog/client"
	"smartlog/client/file"
	"smartlog/client/network"
	"smartlog/uri"
)

func New(u string) (*client.Client, error) {
	ur, err := uri.New(u)
	if err != nil {
		return nil, err
	}
	if ur.Scheme == uri.File {
		return file.New(ur)
	}
	return network.New(ur)
}
