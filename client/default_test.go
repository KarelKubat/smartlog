package client

import (
	"testing"
)

func TestDefaultClient(t *testing.T) {
	for desc, f := range map[string]func(string) error{
		"Info": DefaultClient.Info,
		"Warn": DefaultClient.Warn,
	} {
		if err := f("hello world"); err != nil {
			t.Errorf("DefaultClient.%v(_) = %v, want nil error", desc, err)
		}
	}
}
