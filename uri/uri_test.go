package uri

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		u         string
		wantError string
	}{
		// General format
		{
			u:         "nonsense",
			wantError: "scheme://rest",
		},

		// None needs 1 part
		{
			u:         "none://",
			wantError: "has 0 colon-separated part(s)",
		},
		{
			u:         "none://a:b",
			wantError: "has 2 colon-separated part(s)",
		},

		// File needs 1 part
		{
			u:         "file://",
			wantError: "has 0 colon-separated part(s)",
		},
		{
			u:         "file://a:b",
			wantError: "has 2 colon-separated part(s)",
		},

		// UDP needs 2 parts and a valid port
		{
			u:         "udp://",
			wantError: "has 0 colon-separated part(s)",
		},
		{
			u:         "udp://a:b:c",
			wantError: "has 3 colon-separated part(s)",
		},
		{
			u:         "udp://a:b",
			wantError: "has an invalid port",
		},

		// TCP needs 2 parts and a valid port
		{
			u:         "tcp://",
			wantError: "has 0 colon-separated part(s)",
		},
		{
			u:         "tcp://a:b:c",
			wantError: "has 3 colon-separated part(s)",
		},
		{
			u:         "tcp://a:b",
			wantError: "has an invalid port",
		},
	} {
		_, err := New(test.u)
		if err == nil {
			t.Errorf("New(%q) = _,nil, want something with %q", test.u, test.wantError)
		} else if !strings.Contains(err.Error(), test.wantError) {
			t.Errorf("New(%q) = _,%v, want something with %q", test.u, err, test.wantError)
		}
	}
}

func TestRoundtrip(t *testing.T) {
	for _, u := range []string{
		// Valid URIs
		"none://blackhole",
		"file://stdout",
		"file:///tmp/program.log",
		"udp://:1234",
		"udp://hostname:1234",
		"tcp://:1234",
		"tcp://hostname:1234",
	} {
		ur, err := New(u)
		if err != nil {
			t.Errorf("New(%q) = _,%v, want no error", u, err)
		} else if ur.String() != u {
			t.Errorf("URI from %q stringifies to %q, not identical", u, ur.String())
		}
	}
}
