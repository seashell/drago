package handler

import (
	"net/http"
)

// StatusHandlerAdapter is used to check on server status
type StatusHandlerAdapter struct{}

// NewStatusHandlerAdapter :
func NewStatusHandlerAdapter() *StatusHandlerAdapter {
	return &StatusHandlerAdapter{}
}

// Handle :
func (a *StatusHandlerAdapter) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case http.MethodGet:
		return a.handleGet(rw, req)
	default:
		return nil, NewError(405, ErrMethodNotAllowed)
	}
}

func (a *StatusHandlerAdapter) handleGet(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	return nil, nil
}
