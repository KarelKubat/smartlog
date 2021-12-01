package std

import (
	"os"

	"smartlog/client"
)

func New() *client.Client {
	return &client.Client{
		Type:       client.Std,
		TimeFormat: client.DefaultTimeFormat,
		Writer:     os.Stdout,
	}
}

func ToStd() {
	client.DefaultClient = New()
}

func init() {
	ToStd()
}
