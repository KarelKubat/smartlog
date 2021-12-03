package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"smartlog/uri"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"

	infoTag   byte = 'I'
	warnTag        = 'W'
	errorTag       = 'E'
	separator      = '|'
	space          = ' '
)

type Client struct {
	// Present in all loggers
	TimeFormat string    // defaults to YYYY-MM-DD HH:MM:SS localtime
	Writer     io.Writer // writer for Info(f), Warn(f), Error(f)
	URI        *uri.URI  // URI from which the client was constructed

	Conn       net.Conn // Only in network loggers
	IsTrueFile bool     // Only in file loggers
}

func (c *Client) String() string {
	return fmt.Sprintf("%v", c.URI)
}

func (c *Client) Info(msg string) error {
	return c.sendToWriter(infoTag, msg)
}

func (c *Client) Infof(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Info(full)
}

func (c *Client) Warn(msg string) error {
	return c.sendToWriter(warnTag, msg)
}

func (c *Client) Warnf(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Warn(full)
}

func (c *Client) Error(msg string) error {
	if err := c.sendToWriter(errorTag, msg); err != nil {
		return err
	}
	os.Exit(1)
	return nil // to satisfy the prototype
}

func (c *Client) Errorf(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Error(full)
}

func (c *Client) Passthru(buf []byte) error {
	if c.URI.Scheme == uri.None {
		return nil
	}
	return c.write(buf)
}

func (c *Client) OpenFile() error {
	if c.URI.Scheme != uri.File || c.URI.Parts[0] == "stdout" {
		return fmt.Errorf("%v: attempt to open a file but this URI is not file-based", c)
	}
	var err error
	c.Writer, err = os.OpenFile(c.URI.Parts[0], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("%v: failed to create file: %v", c, err)
	}
	c.IsTrueFile = true
	return nil
}

func (c *Client) timeStamp() []byte {
	now := time.Now()
	var stamp string
	if c.TimeFormat == "" {
		stamp = now.String()
	} else {
		stamp = now.Format(c.TimeFormat)
	}
	return []byte(stamp)
}

func (c *Client) sendToWriter(lev byte, msg string) error {
	if c.URI.Scheme == uri.None {
		return nil
	}

	// Prefix for each line
	prefix := append(c.timeStamp(), space, separator, space, lev, space, separator, space)

	for _, line := range strings.Split(msg, "\n") {
		if line == "" {
			continue
		}
		out := append(append(prefix, []byte(line)...), '\n')
		if err := c.write(out); err != nil {
			return err
		}
	}

	// If the file disappears, reopen it
	if c.IsTrueFile {
		_, err := os.Stat(c.URI.Parts[0])
		if err != nil {
			return c.OpenFile()
		}
	}

	return nil
}

func (c *Client) write(buf []byte) error {
	nWritten := 0
	for nWritten < len(buf) {
		n, err := c.Writer.Write(buf)
		if err != nil {
			return fmt.Errorf("%v: write failure: %v", c, err)
		}
		nWritten += n
	}
	return nil
}
