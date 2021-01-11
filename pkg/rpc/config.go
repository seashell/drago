package rpc

import (
	log "github.com/seashell/drago/pkg/log"
)

type ServerConfig struct {
	//BindAddress
	BindAddress string

	// Logger
	Logger log.Logger

	// Receivers
	Receivers map[string]interface{}
}

func DefaultConfig() *ServerConfig {
	return &ServerConfig{
		BindAddress: "0.0.0.0:8081",
		Receivers:   map[string]interface{}{},
	}
}

func (s *ServerConfig) Merge(b *ServerConfig) *ServerConfig {
	result := *s
	if b.BindAddress != "" {
		result.BindAddress = b.BindAddress
	}
	if b.Logger != nil {
		result.Logger = b.Logger
	}
	if b.Receivers != nil {
		result.Receivers = b.Receivers
	}
	return &result
}

type ClientConfig struct {
	// Logger
	Logger log.Logger

	// URL of the Drago server (e.g. http://127.0.0.1:8081).
	Address string
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Address: "0.0.0.0:8081",
	}
}

func (c *ClientConfig) Merge(b *ClientConfig) *ClientConfig {
	result := *c
	if b.Address != "" {
		result.Address = b.Address
	}
	if b.Logger != nil {
		result.Logger = b.Logger
	}
	return &result
}

// Call is used to make an RPC call to the server
func (c *Client) Call(method string, args interface{}, reply interface{}) error {
	return c.client.Call(method, args, reply)
}
