package any

import (
	"testing"
)

func TestNew(t *testing.T) {
	// This only tests whether New() bubbles up errors from uri.New(). There is not much else to test.
	for _, test := range []struct {
		u         string
		wantError bool
	}{
		{
			u:         "none://whatever",
			wantError: false,
		},
		{
			u:         "nonsense",
			wantError: true,
		},
	} {
		_, err := New(test.u)
		gotError := err != nil
		if gotError != test.wantError {
			t.Errorf("New(%q) = _,error=%v, want error=%v", test.u, gotError, test.wantError)
		}
	}
}
