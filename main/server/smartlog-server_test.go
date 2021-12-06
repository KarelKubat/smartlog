package main

import (
	"os"
	"strings"
	"testing"
)

// FWIW, but it's helpful for myself :-| Usage info with tabs sucks.
func TestUsage(t *testing.T) {
	if strings.Contains(usage, "\t") {
		t.Errorf("usage informatin contains tabs, which leads to ugly output (has your editor forsaken you again)")
	}
}

// We need at least 2 args after the program name. This won't test server startup and clients addition
// but it doesn't need to, that's handled in the modules.
func TestArgs(t *testing.T) {
	errTag := "not enough arguments"
	os.Args = []string{"proggy", "arg"}
	err := run()
	if err == nil {
		t.Errorf("run() with args %v == nil, want error", os.Args)
	} else if !strings.Contains(err.Error(), errTag) {
		t.Errorf("run() with args %v = %v, want something with '%v'", os.Args, err, errTag)
	}
}
