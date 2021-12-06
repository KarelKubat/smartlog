package msg

import (
	"bytes"
	"strings"
	"time"
)

const (
	debugTag   byte = 'D' // Tags for Debug(f) etc.
	infoTag         = 'I'
	warnTag         = 'W'
	fatalTag        = 'F'
	unknownTag      = '?' // For reparsing misses

	separator = '|' // Message parts are separated by space, separator, space
	space     = ' '
)

var (
	DefaultTimeFormat = "2006-01-02 15:04:05 MST"
	UTCTime           = false
)

type MsgType int

const (
	Debug   MsgType = iota // debug level: keep at zero for the tests (or fix msg_test.go)
	Info                   // info level
	Warn                   // warn: won't be dropped
	Fatal                  // terminal: won't be dropped and kills
	Unknown                // keep as last (sentinel) for the tests
)

var tagForType = map[MsgType]byte{
	Debug:   debugTag,
	Info:    infoTag,
	Warn:    warnTag,
	Fatal:   fatalTag,
	Unknown: unknownTag,
}

var typeForTag = map[byte]MsgType{
	debugTag:   Debug,
	infoTag:    Info,
	warnTag:    Warn,
	fatalTag:   Fatal,
	unknownTag: Unknown,
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
		now := time.Now()
		if UTCTime {
			now = now.UTC()
		}
		timestamp = []byte(now.Format(timeFormat))
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

func TypeFromBytes(msg []byte) MsgType {
	parts := bytes.Split(msg, []byte{space, separator, space})
	if len(parts) < 3 || len(parts[1]) != 1 {
		return Unknown
	}
	return typeForTag[parts[1][0]]
}
