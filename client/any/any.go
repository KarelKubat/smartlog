package any

import (
	"errors"
	"smartlog/client"
	"smartlog/client/file"
	"smartlog/client/http"
	"smartlog/client/network"
	"smartlog/client/none"
	"smartlog/uri"
)

func New(u string) (*client.Client, error) {
	ur, err := uri.New(u)
	if err != nil {
		return nil, err
	}
	switch ur.Scheme {
	case uri.None:
		return none.New(ur)
	case uri.File:
		return file.New(ur)
	case uri.TCP:
		return network.New(ur)
	case uri.UDP:
		return network.New(ur)
	case uri.HTTP:
		return http.New(ur)
	}
	return nil, errors.New("internal foobar, unhandled case in any.New")
}
