package any

import (
	"errors"
	"io/fs"
	"strconv"

	"smartlog/client"
	"smartlog/client/file"
	"smartlog/client/network"
	"smartlog/client/std"
	"smartlog/uri"
)

func New(u string) (*client.Client, error) {
	ur, err := uri.New(u)
	if err != nil {
		return nil, err
	}
	switch ur.Scheme {
	case uri.File:
		if ur.Parts[0] == "stdout" {
			return std.New(), nil
		}
		opts := &file.FileClientOpts{
			Filename: ur.Parts[0],
		}
		if len(ur.Parts) > 1 && ur.Parts[1] == "truncate" {
			opts.Truncate = true
		}
		if len(ur.Parts) > 2 {
			umask, err := strconv.ParseInt(ur.Parts[2], 8, 16)
			if err != nil {
				return nil, uri.Error(u, "incorrect umask (must be octal number)")
			}
			opts.UMask = fs.FileMode(umask)
		}
		return file.New(opts)

	case uri.UDP:
		port, err := uri.Port(u, ur.Parts[1])
		if err != nil {
			return nil, err
		}
		return network.New("udp", ur.Parts[0], port)

	case uri.TCP:
		port, err := uri.Port(u, ur.Parts[1])
		if err != nil {
			return nil, err
		}
		return network.New("tcp", ur.Parts[0], port)
	}

	return nil, errors.New("internal foobar, unhandled case in any.New")
}
