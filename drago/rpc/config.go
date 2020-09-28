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
		BindAddress: "0.0.0.0:8181",
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
