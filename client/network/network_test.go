package network

import (
	"testing"

	"github.com/KarelKubat/smartlog/uri"
)

// Test that connection errors are bubbled up. Not much else that we can do.
func TestNew(t *testing.T) {
	_, err := New(&uri.URI{
		Scheme: uri.TCP,
		Parts:  []string{"non.existent.domain", "12345"},
	})
	if err == nil {
		t.Error("New() for nonsense domain = nil, want error")
	}
}
