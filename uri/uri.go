package uri

import (
	"fmt"
	"strconv"
	"strings"
)

// Supported schemata
type URISchema int

const (
	File URISchema = iota
	UDP
	TCP
)

func (u URISchema) String() string {
	return []string{"file", "udp", "tcp"}[u]
}

type URI struct {
	Scheme URISchema
	Parts  []string
}

func New(s string) (*URI, error) {
	uri := &URI{}
	// SCHEME://part1:part2:part3:etc
	top := strings.Split(s, "://")
	if len(top) != 2 {
		return nil, Error(s, "expected: scheme://rest")
	}
	schemeMap := map[string]struct {
		uriType     URISchema
		minParts    int
		maxParts    int
		description string
	}{
		"file": {
			uriType:     File,
			minParts:    1,
			maxParts:    3,
			description: "file://FILENAME or file://FILENAME:truncate or file:://FILENAME:[truncate]:UMASK",
		},
		"udp": {
			uriType:     UDP,
			minParts:    2,
			maxParts:    2,
			description: "udp://SERVER:PORT",
		},
		"tcp": {
			uriType:     TCP,
			minParts:    2,
			maxParts:    2,
			description: "tcp://SERVER:PORT",
		},
	}

	var ok bool
	valid, ok := schemeMap[top[0]]
	if !ok {
		supported := ""
		for key := range schemeMap {
			if supported != "" {
				supported += ","
			}
			supported += key
		}
		return nil, Errorf(s, "unsupported scheme %q, supported: %v", top[0], supported)
	}
	uri.Scheme = valid.uriType
	uri.Parts = strings.Split(top[1], ":")
	if len(uri.Parts) < valid.minParts || len(uri.Parts) > valid.maxParts {
		return nil, Errorf(s, "it has %v colon-separated parts, supported: %v", len(uri.Parts), valid.description)
	}
	return uri, nil
}

func Error(uri, msg string) error {
	return fmt.Errorf("invalid URI %q, %s", uri, msg)
}

func Errorf(uri, format string, args ...interface{}) error {
	return Error(uri, fmt.Sprintf(format, args...))
}

func Port(u, p string) (int, error) {
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, Errorf(u, "port %q is not a number: %v", p, err)
	}
	return port, nil
}
