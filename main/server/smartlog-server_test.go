package main

import (
	"strings"
	"testing"
)

func TestUsage(t *testing.T) {
	// FWIW, but it's helpful for myself :-|
	if strings.Contains(usage, "\t") {
		t.Errorf("usage informatin contains tabs, which leads to ugly output (has your editor forsaken you again)")
	}
}
