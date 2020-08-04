package rpc

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"time"

	log "github.com/seashell/drago/pkg/log"
)

// Server :
type Server struct {
	rpcServer  *rpc.Server
	config     *ServerConfig
	logger     log.Logger
	listener   net.Listener
	listenerCh chan struct{}
}

// NewServer :
func NewServer(config *ServerConfig) (*Server, error) {

	config = DefaultConfig().Merge(config)

	listener, err := net.Listen("tcp", config.BindAddress)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error starting rpc listener: %v", err))
	}

	server := &Server{
		rpcServer:  rpc.NewServer(),
		config:     config,
		listener:   listener,
		logger:     config.Logger,
		listenerCh: make(chan struct{}),
	}

	for name, receiver := range config.Receivers {
		server.rpcServer.RegisterName(name, receiver)
	}

	server.rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	return server, nil
}

// Run :
func (s *Server) Run() {

	fmt.Println(s.listener)
	httpServer := http.Server{
		Addr:              s.listener.Addr().String(),
		Handler:           s.rpcServer,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		defer close(s.listenerCh)
		httpServer.Serve(s.listener)
	}()

	s.logger.Debugf("rpc server started at %s", httpServer.Addr)
}
