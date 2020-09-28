package api

import (
	"net/url"
	"time"
)

const (
	// DefaultAddress is the default Drago server address.
	DefaultAddress = "http://127.0.0.1:8080"

	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 2 * time.Second
)

// Config contains configurations for Drago's API client.
type Config struct {
	// URL of the Drago server (e.g. http://127.0.0.1:8080).
	Address string

	// Token to be used for authentication.
	Token string

	// Request timeout.
	Timeout time.Duration
}

// DefaultConfig returns a default configuration for Drago's API client.
func DefaultConfig() *Config {
	config := &Config{
		Address: DefaultAddress,
	}
	return config
}

// Validate validates the configurations contained wihin the Config struct.
func (c *Config) Validate() error {
	if _, err := url.Parse(c.Address); err != nil {
		return err
	}
	return nil
}

// Merge merges two API client configurations.
func (c *Config) Merge(b *Config) *Config {

	if b == nil {
		return c
	}

	result := *c

	if b.Address != "" {
		result.Address = b.Address
	}
	if b.Token != "" {
		result.Token = b.Token
	}
	if b.Timeout != 0 {
		result.Timeout = b.Timeout
	}

	return &result
}
