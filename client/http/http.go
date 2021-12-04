package http

import (
	h "net/http"
	"strings"

	"smartlog/client"
	"smartlog/uri"
)

const (
	keepMessages = 1024 // # of messages to keep for viewing
)

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
	if len(b.client.Buffer) >= keepMessages {
		b.client.Buffer = b.client.Buffer[:keepMessages-1]
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
