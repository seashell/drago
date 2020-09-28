package client

import (
	"time"

	log "github.com/seashell/drago/pkg/log"
	version "github.com/seashell/drago/version"
)

const (
	defaultLogLevel        = "DEBUG"
	defaultStateDir        = "/etc/drago"
	defaultWireguardPath   = "/home/wireguard"
	defaultInterfacePrefix = "drago"
)

// Config : Drago client configuration
type Config struct {

	// Version is the version of the Drago client.
	Version *version.VersionInfo

	// LogLevel is the level at which the client should output logs.
	LogLevel string

	//Logger is the logger the client will use to log.
	Logger log.Logger

	// Servers is a list of known server addresses. These are as "host:port".
	Servers []string

	// Token contains the auth token used by the client.
	Token string

	// StateDir is the directory to store our state in.
	StateDir string

	// InterfacePrefix is the string prepended to the name of all WireGuard
	// interfaces created by Drago.
	InterfacePrefix string

	// ReconcileInterval is the interval between two reconciliation cycles.
	ReconcileInterval time.Duration

	// WireguardPath is path to the WireGuard binary.
	WireguardPath string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Version:           version.GetVersion(),
		LogLevel:          defaultLogLevel,
		Servers:           []string{"localhost"},
		StateDir:          defaultStateDir,
		InterfacePrefix:   defaultInterfacePrefix,
		ReconcileInterval: 5 * time.Second,
		WireguardPath:     defaultWireguardPath,
	}
}

// Merge combines two config structs, returning the result.
func (c *Config) Merge(b *Config) *Config {
	result := *c

	if b.Logger != nil {
		result.Logger = b.Logger
	}
	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
	}
	if b.Servers != nil {
		result.Servers = b.Servers
	}
	if b.StateDir != "" {
		result.StateDir = b.StateDir
	}
	if b.InterfacePrefix != "" {
		result.InterfacePrefix = b.InterfacePrefix
	}
	if b.ReconcileInterval != 0 {
		result.ReconcileInterval = b.ReconcileInterval
	}
	if b.WireguardPath != "" {
		result.WireguardPath = b.WireguardPath
	}

	return &result
}
