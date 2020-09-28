package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
	structs "github.com/seashell/drago/drago/application/structs"
)

// HandlerFunc is a custom HTTP handler function that returns a struct
// and an error that will be encoded and returned to the client
type HandlerFunc func(Response, *Request) (interface{}, error)

// Middleware
type Middleware func(HandlerFunc) HandlerFunc

// HandlerAdapter
type HandlerAdapter interface {
	http.Handler
}

// Request is a wrapper that extends http.Request with convenience attributes
type Request struct {
	*http.Request

	// Params contains a map of the URL parameters, with keys obtained from the route path
	Params map[string]string
}

// Response is simply a wrapper around the http.ResponseWritter func
type Response struct {
	http.ResponseWriter
}

// BaseHandlerAdapter
type BaseHandlerAdapter struct {
	httprouter.Router
}

// RegisterHandlerFunc
func (a *BaseHandlerAdapter) RegisterHandlerFunc(method string, path string, handler HandlerFunc, middleware ...Middleware) {

	f := func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		params := map[string]string{}
		for _, p := range ps {
			params[p.Key] = p.Value
		}

		// Apply basic middleware to our handler fcn
		fcn := handler

		// Apply custom middleware
		for _, m := range middleware {
			fcn = m(fcn)
		}

		// Wrap net/http structs for more convenient handling
		_resp := Response{rw}
		_req := Request{Request: req, Params: params}

		// Invoke custom handler
		out, err := fcn(_resp, &_req)

		if err != nil {
			code := http.StatusInternalServerError

			if err, ok := err.(Error); ok {
				code = err.Code()
				encoded, err := json.Marshal(&structs.ErrorOutput{
					Message: err.Error(),
				})
				if err != nil {
					fmt.Printf("error encoding json")
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(code)
				rw.Write(encoded)
				return
			}
		}

		encoded, err := json.Marshal(out)
		if err != nil {
			fmt.Printf("error encoding json")
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(encoded)
	}

	a.Handle(method, path, f)
}
