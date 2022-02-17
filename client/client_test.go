package client

import (
	"bytes"
	"strings"
	"testing"

	"github.com/KarelKubat/smartlog/uri"
)

func TestDebugLevels(t *testing.T) {
	for lev := 0; lev < 10; lev++ {
		buf := new(bytes.Buffer)
		cl := &Client{
			URI: &uri.URI{
				Scheme: uri.File,
				Parts:  []string{"buffer"},
			},
			DebugThreshold: uint8(5),
			Writer:         buf,
		}
		if err := cl.Debug(uint8(lev), "hello world"); err != nil {
			t.Fatalf("cl.Debug(%v,_) = %v, need nil error", lev, err)
		}
		if lev > 5 && buf.Len() > 0 {
			t.Errorf("cl.Debug(%v,_) gives output, want none with levels above 5", lev)
		}
	}
}

func TestOpenFile(t *testing.T) {
	for _, test := range []struct {
		filename  string
		wantError string
	}{
		{
			// pseudo-file, can't be opened
			filename:  "stdout",
			wantError: "URI is not file-based",
		},
		{
			// impossible filename
			filename:  "/this/does/not/exist",
			wantError: "failed to create file",
		},
	} {
		cl := &Client{
			URI: &uri.URI{
				Scheme: uri.File,
				Parts:  []string{test.filename},
			},
		}
		err := cl.OpenFile()
		if !strings.Contains(err.Error(), test.wantError) {
			t.Errorf("cl.OpenFile for filename %q = %v, want error with %q", test.filename, err, test.wantError)
		}
	}
}

func TestConnect(t *testing.T) {
	// Don't make this test too long, just 1 attempt and no waittime.
	RestartAttempts = 1
	RestartWait = 0

	for _, test := range []struct {
		uriScheme uri.URISchema
		wantError string
	}{
		{
			uriScheme: uri.None,
			wantError: "internal foobar",
		},
		{
			uriScheme: uri.File,
			wantError: "internal foobar",
		},
		{
			uriScheme: uri.UDP,
			wantError: "failed to (re)connect",
		},
		{
			uriScheme: uri.TCP,
			wantError: "failed to (re)connect",
		},
		{
			uriScheme: uri.HTTP,
			wantError: "internal foobar",
		},
	} {
		cl := &Client{
			URI: &uri.URI{
				Scheme: test.uriScheme,
				Parts:  []string{"non-existent-hostname", "12345"},
			},
		}
		err := cl.Connect()
		if !strings.Contains(err.Error(), test.wantError) {
			t.Errorf("cl.Connect for type %v = %v, want error with %q", test.uriScheme, err, test.wantError)
		}
	}
}
