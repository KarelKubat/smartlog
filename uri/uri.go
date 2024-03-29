package uri

import (
	"fmt"
	"strconv"
	"strings"
)

// Supported schemata
type URISchema int

const (
	None URISchema = iota
	File
	UDP
	TCP
	HTTP
)

func (u URISchema) String() string {
	return []string{"none", "file", "udp", "tcp", "http"}[u]
}

type URI struct {
	Scheme URISchema
	Parts  []string
}

func New(s string) (*URI, error) {
	uri := &URI{}
	// SCHEME://part1:part2:part3:etc, though beyond SCHEME:// we only suport 1 or 2 parts
	top := strings.Split(s, "://")
	if len(top) != 2 {
		return nil, fmt.Errorf("%v: expected: scheme://rest", s)
	}
	schemeMap := map[string]struct {
		uriType     URISchema
		parts       int
		description string
	}{
		"none": {
			uriType:     None,
			parts:       1,
			description: "none://WHATEVER",
		},
		"file": {
			uriType:     File,
			parts:       1,
			description: "file://FILENAME",
		},
		"udp": {
			uriType:     UDP,
			parts:       2,
			description: "udp://SERVER:PORT",
		},
		"tcp": {
			uriType:     TCP,
			parts:       2,
			description: "tcp://SERVER:PORT",
		},
		"http": {
			uriType:     HTTP,
			parts:       2,
			description: "http://SERVER:PORT",
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
			supported += fmt.Sprintf("%s://...", key)
		}
		return nil, fmt.Errorf("%v has an unsupported scheme %q, supported: %v", s, top[0], supported)
	}
	uri.Scheme = valid.uriType
	uri.Parts = strings.Split(top[1], ":")
	nParts := len(uri.Parts)
	if nParts > 0 && uri.Parts[nParts-1] == "" {
		nParts--
	}
	if nParts != valid.parts {
		return nil, fmt.Errorf("%v has %v colon-separated part(s), supported: %v", s, nParts, valid.description)
	}
	if len(uri.Parts) > 1 {
		_, err := strconv.Atoi(uri.Parts[1])
		if err != nil {
			return nil, fmt.Errorf("%v has an invalid port %q (not a number: %v)", s, uri.Parts[1], err)
		}
	}
	return uri, nil
}

func (u *URI) String() string {
	return fmt.Sprintf("%v://%v", u.Scheme, strings.Join(u.Parts, ":"))
}
