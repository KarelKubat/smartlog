package file

import (
	"io/fs"
	"os"

	"smartlog/client"
)

type FileClientOpts struct {
	Filename string
	Truncate bool
	UMask    fs.FileMode
}

func New(opts *FileClientOpts) (*client.Client, error) {
	c := &client.Client{
		Type:       client.File,
		TimeFormat: client.DefaultTimeFormat,
		Filename:   opts.Filename,
	}
	bitmask := os.O_CREATE | os.O_WRONLY
	if opts.Truncate {
		bitmask |= os.O_TRUNC
	} else {
		bitmask |= os.O_APPEND
	}
	umask := opts.UMask
	if umask == 0 {
		umask = 0644
	}
	var err error
	c.Writer, err = os.OpenFile(opts.Filename, bitmask, umask)
	return c, err
}

func ToFile(opts *FileClientOpts) error {
	var err error
	client.DefaultClient, err = New(opts)
	return err
}
