package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	// Default Drago server address
	DefaultAddress = "http://127.0.0.1:8080"
	// Default request timeout
	DefaultTimeout = 2 * time.Second
)

type Config struct {
	// URL of the Drago server (e.g. http://127.0.0.1:8080)
	Address string
	// Token to be used for authentication
	Token string
	// Request timeout
	Timeout time.Duration
}

// Client provides a client to the Drago API
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient returns a new Drago API client
func NewClient(config *Config) (*Client, error) {

	defaults := DefaultConfig()

	if config.Address == "" {
		config.Address = defaults.Address
	} else if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
	}

	client := &Client{
		config:     *config,
		httpClient: cleanhttp.DefaultClient(),
	}

	return client, nil
}

// DefaultConfig returns a default configuration for the client
func DefaultConfig() *Config {
	c := &Config{
		Address: DefaultAddress,
	}
	return c
}
