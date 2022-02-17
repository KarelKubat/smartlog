package http

import (
	"bytes"
	"testing"

	"github.com/KarelKubat/smartlog/client"
)

func TestWrite(t *testing.T) {
	KeepMessages = 10

	for _, test := range []struct {
		nWrites   int
		wantLen   int
		wantFirst int
		wantLast  int
	}{
		{
			// Nothing ongoing
			nWrites:   0,
			wantLen:   0,
			wantFirst: -1, // marker for "not there", see gotFirst/Last below
			wantLast:  -1,
		},
		{
			// 1 write: entry is the first and the last only one
			nWrites:   1,
			wantLen:   1,
			wantFirst: 1,
			wantLast:  1,
		},
		{
			// 9 writes, just less than KeepMessages.
			nWrites:   9,
			wantLen:   9,
			wantFirst: 1,
			wantLast:  9,
		},
		{
			// exactly 10, KeepMessages size
			nWrites:   10,
			wantLen:   10,
			wantFirst: 1,
			wantLast:  10,
		},
		{
			// 11: rollover once
			nWrites:   11,
			wantLen:   10,
			wantFirst: 2,
			wantLast:  11,
		},
		{
			// 21: rollover twice
			nWrites:   21,
			wantLen:   10,
			wantFirst: 12,
			wantLast:  21,
		},
	} {
		cl := &client.Client{
			Writer: new(bytes.Buffer),
		}
		bh := &bufferHandler{
			client: cl,
		}
		for i := 1; i <= test.nWrites; i++ {
			bh.Write([]byte{byte(i)})
		}

		gotFirst := func() int {
			if len(bh.client.Buffer) == 0 {
				return -1
			}
			return int(bh.client.Buffer[0][0])
		}
		gotLast := func() int {
			l := len(bh.client.Buffer)
			if l == 0 {
				return -1
			}
			return int(bh.client.Buffer[l-1][0])
		}
		gotAll := func() []int {
			out := []int{}
			for _, b := range bh.client.Buffer {
				out = append(out, int(b[0]))
			}
			return out
		}
		switch {
		case test.wantLen > 0 && len(bh.client.Buffer) == 0:
			t.Errorf("Write %v times = %v: no entries at all in the buffer", gotAll(), test.nWrites)
		case len(bh.client.Buffer) != test.wantLen:
			t.Errorf("Write %v times= %v: want %v entries in the buffer, got %v", test.nWrites, gotAll(), test.wantLen, len(bh.client.Buffer))
		case gotFirst() != test.wantFirst:
			t.Errorf("Write %v times= %v: first entry mismatch: got %v, want %v", test.nWrites, gotAll(), gotFirst(), test.wantFirst)
		case gotLast() != test.wantLast:
			t.Errorf("Write %v times = %v: last entry mismatch: got %v, want %v", test.nWrites, gotAll(), gotLast(), test.wantLast)
		}
	}
}
