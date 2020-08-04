package handler

import (
	stdhttp "net/http"
	"strings"

	http "github.com/seashell/drago/agent/infrastructure/http"
)

type FilesystemHandlerAdapter struct {
	http.BaseHandlerAdapter
	fsHandler http.HandlerAdapter
}

// NewFilesystemHandlerAdapter creates a new handler adapter for delivering
// static files from a filesystem over HTTP
func NewFilesystemHandlerAdapter(fs stdhttp.FileSystem) http.HandlerAdapter {

	a := &FilesystemHandlerAdapter{
		fsHandler: stdhttp.FileServer(fs),
	}

	a.RegisterHandlerFunc("GET", "/*filepath", a.spa)

	return a
}

func (a *FilesystemHandlerAdapter) spa(resp http.Response, req *http.Request) (interface{}, error) {

	if req.URL.Path != "/" && !strings.HasPrefix(req.URL.Path, "/static") {
		req.URL.Path = "/"
	}
	a.fsHandler.ServeHTTP(resp.ResponseWriter, req.Request)

	return nil, nil
}
