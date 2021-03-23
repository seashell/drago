package client

import (
	"time"

	log "github.com/seashell/drago/pkg/log"
	version "github.com/seashell/drago/version"
)

const (
	defaultLogLevel         = "DEBUG"
	defaultStateDir         = "/tmp/drago"
	defaultWireguardPath    = ""
	defaultInterfacesPrefix = "drago-"
)

// Config : Drago client configuration
type Config struct {

	// DevMode indicates whether the client is running in development mode.
	DevMode bool

	// Name is used to specify the name of the client node.
	Name string

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

	// InterfacesPrefix is the string prepended to the name of all WireGuard
	// interfaces created by Drago.
	InterfacesPrefix string

	// ReconcileInterval is the interval between two reconciliation cycles.
	ReconcileInterval time.Duration

	// WireguardPath is path to the WireGuard binary.
	WireguardPath string

	// Meta contains client metadata
	Meta map[string]string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Name:              "",
		Version:           version.GetVersion(),
		LogLevel:          defaultLogLevel,
		Servers:           []string{"localhost:8081"},
		StateDir:          defaultStateDir,
		InterfacesPrefix:  defaultInterfacesPrefix,
		ReconcileInterval: 5 * time.Second,
		WireguardPath:     defaultWireguardPath,
		Meta:              map[string]string{},
	}
}

// Merge combines two config structs, returning the result.
func (c *Config) Merge(b *Config) *Config {
	result := *c

	if b.Name != "" {
		result.Name = b.Name
	}
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
	if b.InterfacesPrefix != "" {
		result.InterfacesPrefix = b.InterfacesPrefix
	}
	if b.ReconcileInterval != 0 {
		result.ReconcileInterval = b.ReconcileInterval
	}
	if b.WireguardPath != "" {
		result.WireguardPath = b.WireguardPath
	}
	if b.Meta != nil {
		result.Meta = b.Meta
	}

	return &result
}
