package http

import (
	"net/http"
)

// FallthroughHandler :
type FallthroughHandler struct {
	redirectTo string
}

// NewFallthroughHandler :
func NewFallthroughHandler(to string) *FallthroughHandler {
	return &FallthroughHandler{redirectTo: to}
}

// Handle :
func (a *FallthroughHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "GET":
		return a.handleGet(rw, req)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (a *FallthroughHandler) handleGet(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.URL.Path == "/" {
		http.Redirect(rw, req, a.redirectTo, 307)
	} else {
		return nil, NewCodedError(404, ErrNotFound)
	}
	return nil, nil
}
