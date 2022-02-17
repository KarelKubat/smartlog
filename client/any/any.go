package any

import (
	"errors"
	"github.com/KarelKubat/smartlog/client"
	"github.com/KarelKubat/smartlog/client/file"
	"github.com/KarelKubat/smartlog/client/http"
	"github.com/KarelKubat/smartlog/client/network"
	"github.com/KarelKubat/smartlog/client/none"
	"github.com/KarelKubat/smartlog/uri"
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
