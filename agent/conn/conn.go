package conn

import (
	"errors"
	"sync"
	"time"

	log "github.com/seashell/drago/pkg/log"
	rpc "github.com/seashell/drago/pkg/rpc"
)

const (
	errNoServers = "no servers"
)

var (
	ErrNoServers = errors.New(errNoServers)
)

type RPCConnection interface {
	Call(method string, args interface{}, reply interface{}) error
}

type rpcConnection struct {
	logger  log.Logger
	client  *rpc.Client
	address string
	sync.Mutex
}

func NewRPCConnection(address string, logger log.Logger) RPCConnection {
	return &rpcConnection{
		address: address,
		logger:  logger,
	}
}

// Call calls a RPC method on a remote server using the clients RPC client, establishing
// the connection if it's being used for the first time, of if it has been disconnected.
func (c *rpcConnection) Call(method string, args interface{}, reply interface{}) error {

	// Get a cached client or create a new one
	client, err := c.getRPCClient()
	if err != nil {
		return ErrNoServers
	}

	if err := client.Call(method, args, reply); err != nil {
		c.client = nil
		return err
	}

	return nil
}

func (c *rpcConnection) getRPCClient() (*rpc.Client, error) {

	c.Lock()
	defer c.Unlock()

	if c.client != nil {
		return c.client, nil
	}

	client, err := rpc.NewClient(&rpc.ClientConfig{
		Logger:      c.logger,
		Address:     c.address,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	c.client = client

	return c.client, nil
}
