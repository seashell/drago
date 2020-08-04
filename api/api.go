package api

import (
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

// Client provides a client to the Drago API
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient returns a new Drago API client
func NewClient(config *Config) (*Client, error) {

	config = DefaultConfig().Merge(config)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		config:     *config,
		httpClient: cleanhttp.DefaultClient(),
	}

	return client, nil
}
