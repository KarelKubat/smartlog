package file

import (
	"os"

	"github.com/KarelKubat/smartlog/client"
	"github.com/KarelKubat/smartlog/uri"
)

func New(u *uri.URI) (*client.Client, error) {
	c := &client.Client{
		URI: u,
	}
	if u.Parts[0] == "stdout" {
		c.Writer = os.Stdout
		return c, nil
	}
	if err := c.OpenFile(); err != nil {
		return nil, err
	}
	return c, nil
}
