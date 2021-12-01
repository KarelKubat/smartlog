package network

import (
	"fmt"
	"net"

	"smartlog/client"
)

func New(transport, hostname string, port int) (*client.Client, error) {
	c := &client.Client{
		Type:       client.Network,
		TimeFormat: client.DefaultTimeFormat,
	}
	var err error
	c.Conn, err = net.Dial(transport, fmt.Sprintf("%v:%v", hostname, port))
	if err != nil {
		return nil, err
	}
	c.Writer = c.Conn
	return c, nil
}

func ToUDP(hostname string, port int) error {
	var err error
	client.DefaultClient, err = New("udp", hostname, port)
	return err
}

func ToTCP(hostname string, port int) error {
	var err error
	client.DefaultClient, err = New("tcp", hostname, port)
	return err
}
