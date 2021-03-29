package main

import (
	"log"
	"net/rpc"
)

// LeasePlugin :
type LeasePlugin struct {
	logger    log.Logger
	rpcServer *rpc.Server
}

// Config :
type Config struct {
}

// NewLeasePlugin : Creates a new lease plugin object parameterized according to the provided configurations.
func NewLeasePlugin(config *Config) (*LeasePlugin, error) {
	p := &LeasePlugin{}
	return p, nil
}

func main() {
	_, err := NewLeasePlugin(&Config{})
	if err != nil {
		panic(err)
	}
}
