package agent

import (
	"os"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/seashell/drago/version"
)

type Config struct {
	// UI defines whether or not Drago's web UI will be served
	// by the agent
	UI bool `hcl:"ui"`

	// DataDir is the directory to store our state in
	DataDir string `hcl:"data_dir"`

	// PluginDir is the directory where plugins are loaded from
	PluginDir string `hcl:"plugin_dir"`

	// BindAddr is the address on which all of Drago's services will
	// be bound. If not specified, this defaults to 127.0.0.1.
	BindAddr string `hcl:"bind_addr"`

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

// Merge merges two agent configurations
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

	return &result
}

type ServerConfig struct {
	// Enabled controls if the agent is a server
	Enabled bool `hcl:"enabled"`

	// DataDir is the directory where the state will be stored
	DataDir string `hcl:"data_dir"`
}

// Merge is used to merge two server configs together
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

type ClientConfig struct {
	// Enabled controls if the agent is a client
	Enabled bool `hcl:"enabled"`

	// Server is the address of a known Drago server in "host:port" format
	Server string `hcl:"server"`

	// StateDir is the directory where the client state will be kep
	StateDir string `hcl:"data_dir"`

	// InterfacesPrefix is the prefix that will be added to any interface
	InterfacesPrefix string `hcl:"interfaces_prefix"`

	// SyncInterval controls how frequently the client synchronizes its state
	SyncIntervalSeconds time.Duration `hcl:"sync_interval"`
}

func (c *ClientConfig) Merge(b *ClientConfig) *ClientConfig {
	result := *c

	if b.Enabled {
		result.Enabled = true
	}
	if b.Server != "" {
		result.Server = b.Server
	}
	if b.StateDir != "" {
		result.StateDir = b.StateDir
	}
	if b.InterfacesPrefix != "" {
		result.InterfacesPrefix = b.InterfacesPrefix
	}
	if b.SyncIntervalSeconds != 0 {
		result.SyncIntervalSeconds = b.SyncIntervalSeconds
	}

	return &result
}

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

func DefaultConfig() *Config {
	return &Config{
		LogLevel: "INFO",
		UI:       true,
		DataDir:  "/tmp",
		BindAddr: "0.0.0.0",
		Ports: &Ports{
			HTTP: 8080,
			RPC:  8081,
		},
		Server: &ServerConfig{
			Enabled: false,
			DataDir: "/tmp",
		},
		Client: &ClientConfig{
			Enabled:             false,
			InterfacesPrefix:    "dg-",
			SyncIntervalSeconds: 5,
		},
		ACL: &ACLConfig{
			Enabled:  false,
			TokenTTL: 30 * time.Second,
		},
		Version: &version.VersionInfo{},
	}
}

func (c *Config) IsValid() bool {
	// TODO
	return true
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
