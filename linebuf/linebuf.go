package linebuf

import (
	"bytes"
)

type Linebuf struct {
	buf   []byte
	index int
}

func New() *Linebuf {
	l := &Linebuf{}
	l.Reset()
	return l
}

func (l *Linebuf) Add(buf []byte, n int) {
	l.buf = append(l.buf, buf[:n]...)
	l.index = bytes.IndexByte(l.buf, '\n')
}

func (l *Linebuf) Reset() {
	l.buf = []byte{}
	l.index = -1
}

func (l *Linebuf) Complete() bool {
	return l.index >= 0
}

func (l *Linebuf) Bytes() []byte {
	return l.buf
}

func (l *Linebuf) Statement() []byte {
	if !l.Complete() {
		return nil
	}
	stmt := l.buf[:l.index+1]
	l.buf = l.buf[l.index+1:]
	l.index = bytes.IndexByte(l.buf, '\n')

	return stmt
}
