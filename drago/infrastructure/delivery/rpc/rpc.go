package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"time"

	log "github.com/seashell/drago/pkg/log"
)

// Server :
type Server struct {
	config     *ServerConfig
	logger     log.Logger
	listener   net.Listener
	listenerCh chan struct{}
	rpcServer  *rpc.Server
}

// NewServer :
func NewServer(config *ServerConfig) (*Server, error) {

	config = DefaultConfig().Merge(config)

	listener, err := net.Listen("tcp", config.BindAddress)
	if err != nil {
		return nil, fmt.Errorf("error starting rpc listener: %v", err)
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

	httpServer := http.Server{
		Addr:              server.listener.Addr().String(),
		Handler:           server.rpcServer,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		defer close(server.listenerCh)
		httpServer.Serve(server.listener)
	}()

	server.logger.Debugf("rpc server started at %s", httpServer.Addr)

	return server, nil
}
