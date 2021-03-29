package http

import (
	"net/http"
	"strings"
)

// SinglePageApplicationHandler :
type SinglePageApplicationHandler struct {
	fsHandler   http.Handler
	placeholder string
}

// NewSinglePageApplicationHandler : Create a new SPA handler that serves static files
// from a http.FileSystem passed as argument. If the latter is nil, serve a placeholder string.
func NewSinglePageApplicationHandler(fs http.FileSystem, s string) *SinglePageApplicationHandler {

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(s))
	})

	if fs != nil {
		handler = http.FileServer(fs).ServeHTTP
	}

	return &SinglePageApplicationHandler{
		placeholder: s,
		fsHandler:   handler,
	}
}

func (h *SinglePageApplicationHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	path := req.URL.Path

	if path == "" || strings.HasPrefix(path, "/static") {
		h.fsHandler.ServeHTTP(rw, req)
	} else if path == "/" {
		h.fsHandler.ServeHTTP(rw, req)
	} else {
		req.URL.Path = "/"
		h.fsHandler.ServeHTTP(rw, req)
	}

	return nil, nil
}
