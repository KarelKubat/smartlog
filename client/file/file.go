package file

import (
	"os"

	"smartlog/client"
	"smartlog/uri"
)

func New(u *uri.URI) (*client.Client, error) {
	c := &client.Client{
		URI:        u,
		TimeFormat: client.DefaultTimeFormat,
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

func init() {
	client.DefaultClient = &client.Client{
		TimeFormat: client.DefaultTimeFormat,
		Writer:     os.Stdout,
		URI: &uri.URI{
			Scheme: uri.File,
			Parts:  []string{"stdout"},
		},
	}
}
