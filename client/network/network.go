package network

import (
	"fmt"
	"net"
	"strings"

	"smartlog/client"
	"smartlog/uri"
)

func New(ur *uri.URI) (*client.Client, error) {
	c := &client.Client{
		URI:        ur,
		TimeFormat: client.DefaultTimeFormat,
	}
	var err error
	c.Conn, err = net.Dial(fmt.Sprintf("%v", ur.Scheme), strings.Join(ur.Parts, ":"))
	if err != nil {
		return nil, fmt.Errorf("%v: failed to connect: %v", c, err)
	}
	c.Writer = c.Conn
	return c, nil
}
