package http

import (
	log "github.com/seashell/drago/pkg/log"
)

// Config contains configurations for the http module
type Config struct {
	// Server bind address in the form host:port
	BindAddress string

	// Handlers contains a mapping of paths to http.HandlerFunc
	Handlers map[string]Handler

	// Handlers contains middleware functions to be applied to handlers
	Middleware []Middleware

	// Logger
	Logger log.Logger
}

// DefaultConfig :
func DefaultConfig() *Config {
	return &Config{
		BindAddress: "0.0.0.0:9876",
		Handlers:    map[string]Handler{},
		Middleware:  []Middleware{},
	}
}

// Merge :
func (s *Config) Merge(b *Config) *Config {
	result := *s
	if b.BindAddress != "" {
		result.BindAddress = b.BindAddress
	}
	if b.Logger != nil {
		result.Logger = b.Logger
	}
	if b.Handlers != nil {
		result.Handlers = b.Handlers
	}
	if b.Middleware != nil {
		result.Middleware = b.Middleware
	}
	return &result
}
