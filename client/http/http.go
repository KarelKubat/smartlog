package http

import (
	h "net/http"
	"strings"

	"github.com/KarelKubat/smartlog/client"
	"github.com/KarelKubat/smartlog/uri"
)

var KeepMessages = 1024 // # of messages to keep for viewing)

func New(ur *uri.URI) (*client.Client, error) {
	c := &client.Client{
		URI:    ur,
		Buffer: [][]byte{},
	}
	wr := &bufferHandler{
		client: c,
	}
	c.Writer = wr

	mux := h.NewServeMux()
	mux.Handle("/", wr)
	go func() {
		h.ListenAndServe(strings.Join(ur.Parts, ":"), mux)
	}()

	return c, nil
}

type bufferHandler struct {
	client *client.Client
}

func (b *bufferHandler) Write(p []byte) (int, error) {
	if len(b.client.Buffer) >= KeepMessages {
		b.client.Buffer = b.client.Buffer[1:KeepMessages]
	}
	b.client.Buffer = append(b.client.Buffer, p)
	return len(p), nil
}

func (b *bufferHandler) ServeHTTP(w h.ResponseWriter, r *h.Request) {
	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte("<pre>"))
	for _, b := range b.client.Buffer {
		w.Write(b)
	}
	w.Write([]byte("</pre>"))
}
