package handler

import (
	"net/http"
)

// FallthroughHandlerAdapter :
type FallthroughHandlerAdapter struct {
	redirectTo string
}

// NewFallthroughHandlerAdapter :
func NewFallthroughHandlerAdapter(to string) *FallthroughHandlerAdapter {
	return &FallthroughHandlerAdapter{redirectTo: to}
}

// Handle :
func (a *FallthroughHandlerAdapter) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case http.MethodGet:
		return a.handleGet(rw, req)
	default:
		return nil, NewError(405, ErrInvalidMethod)
	}
}

func (a *FallthroughHandlerAdapter) handleGet(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.URL.Path == "/" {
		http.Redirect(rw, req, a.redirectTo, 307)
	} else {
		return nil, NewError(404, ErrNotFound)
	}
	return nil, nil
}
