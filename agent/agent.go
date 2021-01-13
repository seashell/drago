package agent

import (
	"errors"
	"fmt"
	"sync"

	client "github.com/seashell/drago/client"
	drago "github.com/seashell/drago/drago"
	config "github.com/seashell/drago/drago/structs/config"
	log "github.com/seashell/drago/pkg/log"
)

// Agent :
type Agent struct {
	config *Config
	logger log.Logger
	server *drago.Server
	client *client.Client

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New creates a new Drago agent from the configuration,
// potentially returning an error
func New(config *Config, logger log.Logger) (*Agent, error) {

	config = DefaultConfig().Merge(config)

	if logger == nil {
		return nil, errors.New("missing logger")
	}

	a := &Agent{
		config:     config,
		logger:     logger.WithName("agent"),
		shutdownCh: make(chan struct{}),
	}

	// Setup Drago server
	if err := a.setupServer(); err != nil {
		return nil, err
	}

	// Setup Drago client
	if err := a.setupClient(); err != nil {
		return nil, err
	}

	// Make sure agent will be running at least as a client or as a server
	if a.client == nil && a.server == nil {
		return nil, errors.New("must have either client or server mode enabled")
	}

	return a, nil
}

// Shutdown is used to terminate the agent.
func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}

	a.logger.Infof("requesting shutdown")
	if a.client != nil {
		if err := a.client.Shutdown(); err != nil {
			a.logger.Errorf("client shutdown failed: %s", err.Error())
		}
	}
	if a.server != nil {
		if err := a.server.Shutdown(); err != nil {
			a.logger.Errorf("server shutdown failed: %s", err.Error())
		}
	}

	a.logger.Infof("agent shutdown complete")

	a.shutdown = true
	close(a.shutdownCh)

	return nil
}

// Setup Drago server, if enabled
func (a *Agent) setupServer() error {

	if !a.config.Server.Enabled {
		return nil
	}

	config, err := a.serverConfig()
	if err != nil {
		return fmt.Errorf("server config setup failed: %v", err)
	}

	server, err := drago.NewServer(config)
	if err != nil {
		return fmt.Errorf("server setup failed: %v", err)
	}

	a.server = server

	return nil
}

// clientConfig creates a new drago.Config struct based on an
// agent.Config struct and which can be used to initialize
// a Drago client
func (a *Agent) serverConfig() (*drago.Config, error) {

	c := drago.DefaultConfig()

	c.UI = a.config.UI
	c.BindAddr = a.config.BindAddr
	c.DataDir = a.config.Server.DataDir

	c.Ports = &drago.Ports{
		HTTP: a.config.Ports.HTTP,
		RPC:  a.config.Ports.RPC,
	}

	c.ACL = &config.ACLConfig{
		Enabled:  a.config.ACL.Enabled,
		TokenTTL: a.config.ACL.TokenTTL,
	}

	c.LogLevel = a.config.LogLevel
	c.Logger = a.logger

	return c, nil
}

// Setup Drago client, if enabled
func (a *Agent) setupClient() error {
	if !a.config.Client.Enabled {
		return nil
	}

	config, err := a.clientConfig()
	if err != nil {
		return fmt.Errorf("client config setup failed: %v", err)
	}

	client, err := client.New(config)
	if err != nil {
		return fmt.Errorf("client setup failed: %v", err)
	}

	a.client = client

	return nil
}

// clientConfig creates a new client.Config struct based on an
// agent.Config struct and which can be used to initialize
// a Drago client
func (a *Agent) clientConfig() (*client.Config, error) {
	c := client.DefaultConfig()

	c.Servers = a.config.Client.Servers
	c.StateDir = a.config.Server.DataDir

	c.LogLevel = a.config.LogLevel
	c.Logger = a.logger

	return c, nil
}
