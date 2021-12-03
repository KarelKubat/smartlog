package msg

import (
	"strings"
	"time"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"

	debugTag  byte = 'D'
	infoTag        = 'I'
	warnTag        = 'W'
	errorTag       = 'E'
	separator      = '|'
	space          = ' '
)

type MsgType int

const (
	Debug MsgType = iota
	Info
	Warn
	Error
)

var tagForType = map[MsgType]byte{
	Debug: debugTag,
	Info:  infoTag,
	Warn:  warnTag,
	Error: errorTag,
}

type Message struct {
	Type       MsgType
	TimeFormat string
	Timestamp  []byte
	Message    string
}

func BytesFromMessage(m *Message) [][]byte {
	timestamp := []byte(m.Timestamp)
	if len(timestamp) == 0 {
		timeFormat := m.TimeFormat
		if timeFormat == "" {
			timeFormat = DefaultTimeFormat
		}
		timestamp = []byte(time.Now().Format(timeFormat))
	}
	prefix := append(timestamp, space, separator, space, tagForType[m.Type], space, separator, space)
	out := [][]byte{}
	for _, line := range strings.Split(m.Message, "\n") {
		if line == "" {
			continue
		}
		lineBytes := append(append(prefix, []byte(line)...), '\n')
		out = append(out, lineBytes)
	}
	return out
}
