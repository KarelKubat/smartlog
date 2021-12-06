package server

import (
	"strings"
	"testing"
)

// We can only test the bubbling up of errors. Intergration tests are handled elsewhere.
func TestNew(t *testing.T) {
	for _, test := range []struct {
		u         string
		wantError string
	}{
		{
			// Errors from uri.New() are bubbled up
			u:         "nonsense",
			wantError: "scheme://rest",
		},
		{
			// Only tcp:// or udp:// are allowed
			u:         "file://stdout",
			wantError: "only udp:// or tcp://",
		},
	} {
		_, err := New(test.u)
		if !strings.Contains(err.Error(), test.wantError) {
			t.Errorf("New(%q) = _,%v, want something with %q", test.u, err, test.wantError)
		}
	}
}
