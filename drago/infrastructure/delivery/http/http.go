package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/seashell/drago/drago/application/structs"
	"github.com/seashell/drago/pkg/log"
)

// Error :
type Error interface {
	error
	Code() int
}

// Server :
type Server struct {
	config     *Config
	logger     log.Logger
	listener   net.Listener
	listenerCh chan struct{}
	mux        *http.ServeMux
}

// NewServer :
func NewServer(config *Config) (*Server, error) {

	config = DefaultConfig().Merge(config)

	listener, err := net.Listen("tcp", config.BindAddress)
	if err != nil {
		return nil, fmt.Errorf("error starting HTTP listener: %v", err)
	}

	server := &Server{
		config:     config,
		listener:   listener,
		logger:     config.Logger,
		listenerCh: make(chan struct{}),
		mux:        http.NewServeMux(),
	}

	for pattern, handler := range config.Handlers {

		fcn := httpHandlerFunc(handler)

		// Apply custom middleware
		for _, m := range server.config.Middleware {
			fcn = m(fcn)
		}

		fcn = http.StripPrefix(strings.TrimSuffix(pattern, "/"), fcn).ServeHTTP

		server.mux.HandleFunc(pattern, fcn)
	}

	httpServer := http.Server{
		Addr:    server.listener.Addr().String(),
		Handler: server.mux,
	}

	go func() {
		defer close(server.listenerCh)
		httpServer.Serve(server.listener)
	}()

	server.logger.Debugf("http server started at %s", httpServer.Addr)

	return server, nil
}

// httpHandlerFunc converts a custom handler func to http.HandlerFunc
func httpHandlerFunc(handler Handler) http.HandlerFunc {

	f := func(rw http.ResponseWriter, req *http.Request) {

		fcn := handler.Handle

		// Invoke custom handler
		out, err := fcn(rw, req)

		if err != nil {
			code := http.StatusInternalServerError

			if err, ok := err.(Error); ok {
				code = err.Code()
				encoded := encode(&structs.Error{
					Message: err.Error(),
				})
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(code)
				rw.Write(encoded)
				return
			}
		}

		if out != nil {
			encoded := encode(out)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(encoded)
		} else {
			rw.WriteHeader(http.StatusNoContent)
		}

	}

	return f
}

func encode(in interface{}) []byte {
	encoded, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return encoded
}
