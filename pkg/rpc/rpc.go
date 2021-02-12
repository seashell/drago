package rpc

import (
	"fmt"
	"net"
	"net/rpc"

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

	// Use tls.Listen for serving with TLS
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

	go func() {
		for {
			// TODO: Handle signals
			select {
			default:
			}
			conn, _ := listener.Accept()
			cdc := NewMsgpackServerCodec(conn)
			go func() {
				server.rpcServer.ServeCodec(cdc)
			}()
		}
	}()

	server.logger.Debugf("rpc server started at %s", config.BindAddress)

	return server, nil
}

// Client :
type Client struct {
	config *ClientConfig
	logger log.Logger
	client *rpc.Client
}

// NewClient :
func NewClient(config *ClientConfig) (*Client, error) {
	config = DefaultClientConfig().Merge(config)

	c := &Client{
		config: config,
		logger: config.Logger,
	}

	// Use tls.Dial for connection with TLS
	conn, err := net.Dial("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	cdc := NewMsgpackClientCodec(conn)
	c.client = rpc.NewClientWithCodec(cdc)

	return c, nil
}
