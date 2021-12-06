package msg

import (
	"bytes"
	"testing"
	"time"
)

// Check that tagforType and typeForTag are complete.
func TestMaps(t *testing.T) {
	for tp := Debug; tp <= Unknown; tp++ {
		tag, ok := tagForType[tp]
		if !ok {
			t.Errorf("tagForType[%v] is not defined", tp)
		}
		tp1, ok := typeForTag[tag]
		if !ok {
			t.Errorf("typeForTag[%v] is not defined", tag)
		}
		if tp != tp1 {
			t.Errorf("tagForType[%v] = %v, typeForTag[%v] = %v, but %v!=%v", tp, tag, tag, tp1, tp, tp1)
		}
	}
}

func TestBytesFromMessage(t *testing.T) {
	// Check that the returned [][]byte corresponds with the # of lines we send in.
	for _, test := range []struct {
		message   string
		wantLines int
	}{
		{
			// oneliner
			message:   "hello world",
			wantLines: 1,
		},
		{
			// 2 lines
			message:   "hello\nworld",
			wantLines: 2,
		},
		{
			// empty lines get skipped
			message:   "\n\n\n\nhello\nworld\n",
			wantLines: 2,
		},
	} {
		b := BytesFromMessage(&Message{
			Type:    Warn,
			Message: test.message,
		})
		if len(b) != test.wantLines {
			t.Errorf("BytesFromMessage for %q = %v lines, want %v", test.message, len(b), test.wantLines)
		}
	}

	// Check that UTCTime is correctly handled.
	msg := &Message{
		Type:       Warn,
		TimeFormat: "",
		Timestamp:  []byte(time.Now().Format("")),
		Message:    "hello world",
	}
	UTCTime = false
	b1 := BytesFromMessage(msg)
	UTCTime = true
	b2 := BytesFromMessage(msg)
	if cmp := bytes.Compare(b1[0], b2[0]); cmp == 0 {
		t.Errorf("UTC time is picked up incorrectly: %v vs %v", string(b1[0]), string(b2[0]))
	}
}

func TestTypeFromBytes(t *testing.T) {
	msg := &Message{
		TimeFormat: "",
		Timestamp:  []byte(time.Now().Format("")),
		Message:    "hello world",
	}
	for tp := Debug; tp <= Unknown; tp++ {
		msg.Type = tp
		b := BytesFromMessage(msg)
		out := TypeFromBytes(b[0])
		if out != tp {
			t.Errorf("TypeFromBytes with input type %v = %v, want equals", out, tp)
		}
	}
}
