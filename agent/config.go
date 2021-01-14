package agent

import (
	"os"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/seashell/drago/version"
)

// Config contains configurations for the Drago agent
type Config struct {
	// UI defines whether or not Drago's web UI will be served
	// by the agent
	UI bool `hcl:"ui"`

	// Name is used to identify individual agents
	Name string `hcl:"name"`

	// DataDir is the directory to store our state in
	DataDir string `hcl:"data_dir"`

	// PluginDir is the directory where plugins are loaded from
	PluginDir string `hcl:"plugin_dir"`

	// BindAddr is the address on which all of Drago's services will
	// be bound. If not specified, this defaults to 127.0.0.1.
	BindAddr string `hcl:"bind_addr"`

	// AdvertiseAddrs is a struct containing the addresses advertised
	// for each of Drago's network services in host:port format.
	// It is optional, and all addresses default to the bind address
	// with the default port corresponding to each service.
	AdvertiseAddrs *AdvertiseAddrs `hcl:"advertise_addrs"`

	// LogLevel is the level of the logs to put out
	LogLevel string `hcl:"log_level"`

	// Ports is used to control the network ports we bind to.
	Ports *Ports `hcl:"ports"`

	// Server contains all server-specific configurations
	Server *ServerConfig `hcl:"server"`

	// Client contains all client-specific configurations
	Client *ClientConfig `hcl:"client"`

	// ACL contains all ACL configurations
	ACL *ACLConfig `hcl:"acl"`

	// DevMode is set by the --dev CLI flag.
	DevMode bool `hcl:"-"`

	// Version information (set at compilation time)
	Version *version.VersionInfo
}

// Merge merges two Config structs, returning the result
func (c *Config) Merge(b *Config) *Config {

	if b == nil {
		return c
	}

	result := *c

	if b.UI {
		result.UI = true
	}
	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
	}
	if b.Name != "" {
		result.Name = b.Name
	}
	if b.DataDir != "" {
		result.DataDir = b.DataDir
	}
	if b.BindAddr != "" {
		result.BindAddr = b.BindAddr
	}
	if b.Version != nil {
		result.Version = b.Version
	}

	// Apply the ports config
	if result.Ports == nil && b.Ports != nil {
		ports := *b.Ports
		result.Ports = &ports
	} else if b.Ports != nil {
		result.Ports = result.Ports.Merge(b.Ports)
	}

	// Apply the client config
	if result.Client == nil && b.Client != nil {
		client := *b.Client
		result.Client = &client
	} else if b.Client != nil {
		result.Client = result.Client.Merge(b.Client)
	}

	// Apply the server config
	if result.Server == nil && b.Server != nil {
		server := *b.Server
		result.Server = &server
	} else if b.Server != nil {
		result.Server = result.Server.Merge(b.Server)
	}

	// Apply the advertise addrs config
	if result.AdvertiseAddrs == nil && b.AdvertiseAddrs != nil {
		advertise := *b.AdvertiseAddrs
		result.AdvertiseAddrs = &advertise
	} else if b.AdvertiseAddrs != nil {
		result.AdvertiseAddrs = result.AdvertiseAddrs.Merge(b.AdvertiseAddrs)
	}

	return &result
}

// ServerConfig contains configurations for the Drago server
type ServerConfig struct {
	// Enabled controls if the agent is a server
	Enabled bool `hcl:"enabled"`

	// DataDir is the directory where the state will be stored
	DataDir string `hcl:"data_dir"`
}

// Merge merges two ServerConfig structs, returning the result
func (c *ServerConfig) Merge(b *ServerConfig) *ServerConfig {
	result := *c

	if b.Enabled {
		result.Enabled = true
	}
	if b.DataDir != "" {
		result.DataDir = b.DataDir
	}
	return &result
}

// ClientConfig contains configurations for the Drago client
type ClientConfig struct {
	// Enabled controls if the agent is a client
	Enabled bool `hcl:"enabled"`

	// Server is the address of a known Drago server in "host:port" format
	Servers []string `hcl:"servers"`

	// StateDir is the directory where the client state will be kep
	StateDir string `hcl:"data_dir"`

	// InterfacesPrefix is the prefix that will be added to all WireGuard
	// interfaces managed by Drago
	InterfacesPrefix string `hcl:"interfaces_prefix"`

	// WireguardPath is the path to a userspace WireGuard binary, if available
	WireguardPath string `hcl:"wireguard_path"`

	// Meta contains metadata about the client node
	Meta map[string]string `hcl:"meta"`

	// SyncInterval controls how frequently the client synchronizes its state
	SyncIntervalSeconds time.Duration `hcl:"sync_interval"`
}

// Merge merges two ClientConfig structs, returning the result
func (c *ClientConfig) Merge(b *ClientConfig) *ClientConfig {
	result := *c

	if b.Enabled {
		result.Enabled = true
	}
	if b.Servers != nil {
		result.Servers = b.Servers
	}
	if b.StateDir != "" {
		result.StateDir = b.StateDir
	}
	if b.WireguardPath != "" {
		result.WireguardPath = b.WireguardPath
	}
	if b.InterfacesPrefix != "" {
		result.InterfacesPrefix = b.InterfacesPrefix
	}
	if b.SyncIntervalSeconds != 0 {
		result.SyncIntervalSeconds = b.SyncIntervalSeconds
	}
	if b.Meta != nil {
		result.Meta = b.Meta
	}

	return &result
}

// ACLConfig contains configuration for Drago's ACL
type ACLConfig struct {
	// Enabled controls if the ACLs are managed and enforced
	Enabled bool `hcl:"enabled"`
	// TokenTTL controls how long we cache ACL tokens. This controls
	// how stale they can be when we are enforcing policies. Defaults
	// to "30s". Reducing this impacts performance by forcing more
	// frequent resolution.
	TokenTTL    time.Duration
	TokenTTLHCL string `hcl:"token_ttl" json:"-"`
}

// Merge merges two ACLConfig structs, returning the result
func (c *ACLConfig) Merge(b *ACLConfig) *ACLConfig {
	result := *c

	if b.Enabled {
		result.Enabled = true
	}
	if b.TokenTTL != 0 {
		result.TokenTTL = b.TokenTTL
	}
	if b.TokenTTLHCL != "" {
		result.TokenTTLHCL = b.TokenTTLHCL
	}
	return &result
}

// Ports encapsulates the various ports we bind to for network services. If any
// are not specified then the defaults are used instead.
type Ports struct {
	HTTP int `hcl:"http"`
	RPC  int `hcl:"rpc"`
}

// Merge merges two Ports structs, returning the result
func (c *Ports) Merge(b *Ports) *Ports {
	result := *c

	if b.HTTP != 0 {
		result.HTTP = b.HTTP
	}
	if b.RPC != 0 {
		result.RPC = b.RPC
	}

	return &result
}

// AdvertiseAddrs is used to control the addresses Drago advertises for
// its different network services. All are optional and default to BindAddr and
// their default Port. Expected format is address:port.
type AdvertiseAddrs struct {
	Peer   string `hcl:"peer"`
	Client string `hcl:"client"`
}

// Merge merges two AdvertiseAddrs structs, returning the result
func (c *AdvertiseAddrs) Merge(b *AdvertiseAddrs) *AdvertiseAddrs {
	result := *c

	if b.Peer != "" {
		result.Peer = b.Peer
	}
	if b.Client != "" {
		result.Client = b.Client
	}

	return &result
}

// DefaultConfig returns a Config struct populated with sane defaults
func DefaultConfig() *Config {
	return &Config{
		LogLevel: "DEBUG",
		UI:       true,
		Name:     "",
		DataDir:  "/tmp/drago",
		BindAddr: "0.0.0.0",
		Ports: &Ports{
			HTTP: 8080,
			RPC:  8081,
		},
		AdvertiseAddrs: &AdvertiseAddrs{},
		Server: &ServerConfig{
			Enabled: false,
			DataDir: "/tmp/drago",
		},
		Client: &ClientConfig{
			Enabled:             false,
			Servers:             []string{"127.0.0.1:8081"},
			Meta:                map[string]string{},
			InterfacesPrefix:    "drago",
			WireguardPath:       "",
			SyncIntervalSeconds: 5,
		},
		ACL: &ACLConfig{
			Enabled:  false,
			TokenTTL: 30 * time.Second,
		},
		Version: version.GetVersion(),
	}
}

// EmptyConfig returns an empty Config struct with all nested structs
// also initialized to a non-nil empty value.
func EmptyConfig() *Config {
	return &Config{
		Ports:          &Ports{},
		AdvertiseAddrs: &AdvertiseAddrs{},
		Server:         &ServerConfig{},
		Client:         &ClientConfig{},
		ACL:            &ACLConfig{},
	}
}

// Validate returns an error in case a Config struct is invalid.
func (c *Config) Validate() error {
	// TODO: implement validation
	return nil
}

// LoadFromFile loads the configuration from a given path
func (c *Config) LoadFromFile(path string) (*Config, error) {

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = hclsimple.DecodeFile(path, nil, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
