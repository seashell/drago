package drago

import (
	"time"
)

const (
	DefaultHttpPort = 8080
)

// Config : Drago server configuration
type Config struct {
	// BindAddr
	BindAddr string

	// DataDir is the directory to store our state in
	DataDir string

	// HostGCInterval is how often we dispatch a job to GC inactive hosts.
	HostGCInterval time.Duration
}

// DefaultConfig returns the default configuration. Only used as the basis for
// merging agent or test parameters.
func DefaultConfig() *Config {
	return &Config{}
}
