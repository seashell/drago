package http

import (
	"net/http"
)

// Handler :
type Handler interface {
	Handle(http.ResponseWriter, *http.Request) (interface{}, error)
}

// HandlerFunc is a custom HTTP handler function that returns a struct
// and an error that will be encoded and returned to the client
type HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, error)
