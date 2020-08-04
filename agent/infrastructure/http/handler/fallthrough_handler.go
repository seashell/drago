package handler

import (
	stdhttp "net/http"

	http "github.com/seashell/drago/agent/infrastructure/http"
)

type FallthroughHandlerAdapter struct {
	http.BaseHandlerAdapter
}

func NewFallthroughHandlerAdapter(to string) http.HandlerAdapter {
	a := &FallthroughHandlerAdapter{}

	a.HandlerFunc("GET", "/", func(rw stdhttp.ResponseWriter, req *stdhttp.Request) {
		if req.URL.Path == "/" {
			stdhttp.Redirect(rw, req, to, stdhttp.StatusTemporaryRedirect)
		} else {
			rw.WriteHeader(stdhttp.StatusNotFound)
		}
	})

	return a
}
