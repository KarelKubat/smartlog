package file

import (
	"io"
	"os"
	"testing"

	"smartlog/uri"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		filename   string
		wantError  bool
		wantWriter io.Writer
	}{
		{
			filename:   "stdout",
			wantError:  false,
			wantWriter: os.Stdout,
		},
		{
			filename:   "/non/existing/file/that/cannot/be/opened",
			wantError:  true,
			wantWriter: nil,
		},
	} {
		cl, err := New(&uri.URI{
			Scheme: uri.File,
			Parts:  []string{test.filename},
		})
		gotError := err != nil
		if gotError != test.wantError {
			t.Errorf("filename %q: New() = _,error=%v, want error=%v", test.filename, gotError, test.wantError)
		}
		if cl != nil && cl.Writer != test.wantWriter {
			t.Errorf("filename %q: New() doesn't yield the expected writer", test.filename)
		}
	}
}
