package playground

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// BrowserRenderer runs the same bounded preview pipeline without listening on a
// socket. readAsset must return only trusted, deployment-owned stylesheet bytes;
// it must never fetch URLs supplied by an authored screen.
type BrowserRenderer struct{ server *server }

func NewBrowserRenderer(readAsset func(string) ([]byte, error)) (*BrowserRenderer, error) {
	if readAsset == nil {
		return nil, fmt.Errorf("trusted stylesheet reader required")
	}
	registry, err := newRegistry()
	if err != nil {
		return nil, err
	}
	return &BrowserRenderer{server: &server{registry: registry, busy: make(chan struct{}, 1), readAsset: readAsset}}, nil
}

// Request accepts the original JSON text so JavaScript cannot round JSON numbers
// before Go's UseNumber decoder validates and renders them. Its result is an HTTP-
// shaped JSON envelope, with the response body retained as a JSON string.
func (b *BrowserRenderer) Request(method, path, body string) string {
	w := &browserResponse{header: make(http.Header)}
	if method != "POST" || path != "/api/preview" {
		respond(w, 404, map[string]string{"error": "Unknown browser renderer endpoint"})
	} else {
		r, err := http.NewRequest(method, path, strings.NewReader(body))
		if err != nil {
			respond(w, 400, map[string]string{"error": err.Error()})
		} else {
			r.Header.Set("Content-Type", "application/json")
			b.server.preview(w, r)
		}
	}
	encoded, _ := json.Marshal(struct {
		Status int    `json:"status"`
		Body   string `json:"body"`
	}{w.status, w.body.String()})
	return string(encoded)
}

type browserResponse struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func (w *browserResponse) Header() http.Header { return w.header }
func (w *browserResponse) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
}
func (w *browserResponse) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(p)
}
