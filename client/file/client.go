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
	var err error
	c.Writer, err = os.OpenFile(u.Parts[0], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	return c, err
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
