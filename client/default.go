package client

import (
	"os"

	"github.com/KarelKubat/smartlog/uri"
)

var DefaultClient *Client

func Debug(lev uint8, msg string) error {
	return DefaultClient.Debug(lev, msg)
}

func Debugf(lev uint8, format string, args ...interface{}) error {
	return DefaultClient.Debugf(lev, format, args...)
}

func Info(msg string) error {
	return DefaultClient.Info(msg)
}

func Infof(format string, args ...interface{}) error {
	return DefaultClient.Infof(format, args...)
}

func Warn(msg string) error {
	return DefaultClient.Warn(msg)
}

func Warnf(format string, args ...interface{}) error {
	return DefaultClient.Warnf(format, args...)
}

func Fatal(msg string) error {
	return DefaultClient.Fatal(msg)
}

func Fatalf(format string, args ...interface{}) error {
	return DefaultClient.Fatalf(format, args...)
}

func init() {
	DefaultClient = &Client{
		Writer: os.Stdout,
		URI: &uri.URI{
			Scheme: uri.File,
			Parts:  []string{"stdout"},
		},
	}
}
