package agent

import (
	"errors"
	"fmt"
	stdhttp "net/http"
	"sync"

	handler "github.com/seashell/drago/agent/adapter/http"
	middleware "github.com/seashell/drago/agent/adapter/http/middleware"
	conn "github.com/seashell/drago/agent/conn"
	client "github.com/seashell/drago/client"
	drago "github.com/seashell/drago/drago"
	config "github.com/seashell/drago/drago/structs/config"
	http "github.com/seashell/drago/pkg/http"
	log "github.com/seashell/drago/pkg/log"
)

// Agent :
type Agent struct {
	config *Config
	logger log.Logger

	rpcConn conn.RPCConnection

	server *drago.Server
	client *client.Client

	httpServer *http.Server

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

	if err := a.setupRPCConnection(); err != nil {
		return nil, err
	}

	// Setup Drago server
	if err := a.setupServer(); err != nil {
		return nil, err
	}

	// Setup Drago client
	if err := a.setupClient(); err != nil {
		if a.server != nil {
			a.server.Shutdown()
		}
		return nil, err
	}

	// Make sure agent will be running at least as a client or as a server
	if a.client == nil && a.server == nil {
		return nil, errors.New("must have either client or server mode enabled")
	}

	if err := a.setupHTTPServer(); err != nil {
		a.Shutdown()
		return nil, fmt.Errorf("could not initialize http server: %s", err)
	}

	return a, nil
}

// Stats returns a map containing relevant stats
func (a *Agent) Stats() map[string]map[string]string {

	stats := map[string]map[string]string{}

	if a.server != nil {
		subStat := a.server.Stats()
		for k, v := range subStat {
			stats[k] = v
		}
	}
	if a.client != nil {
		subStat := a.client.Stats()
		for k, v := range subStat {
			stats[k] = v
		}
	}

	return stats
}

// Config returns a copy of the agent's Config struct
func (a *Agent) Config() map[string]interface{} {
	config := map[string]interface{}{}
	// TODO populate config map
	return config
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

func (a *Agent) setupRPCConnection() error {

	address := ""

	if a.config.Server.Enabled {
		address = fmt.Sprintf("%s:%d", a.config.BindAddr, a.config.Ports.RPC)
	} else {
		address = a.config.Client.Servers[0]
	}

	a.rpcConn = conn.NewRPCConnection(address, a.logger)

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

// serverConfig creates a new drago.Config struct based on an
// agent.Config struct and which can be used to initialize
// a Drago server
func (a *Agent) serverConfig() (*drago.Config, error) {

	c := drago.DefaultConfig()

	c.UI = a.config.UI
	c.DevMode = *a.config.DevMode
	c.BindAddr = a.config.BindAddr
	c.DataDir = a.config.DataDir

	c.Ports = &drago.Ports{
		HTTP: a.config.Ports.HTTP,
		RPC:  a.config.Ports.RPC,
	}

	c.ACL = &config.ACLConfig{
		Enabled: a.config.ACL.Enabled,
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

	client, err := client.New(a.rpcConn, config)
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

	c.Name = a.config.Name

	c.Servers = a.config.Client.Servers
	c.StateDir = a.config.Client.StateDir

	c.AdvertiseAddress = a.config.AdvertiseAddrs.Peer

	c.WireguardPath = a.config.Client.WireguardPath
	c.InterfacesPrefix = a.config.Client.InterfacesPrefix

	c.Meta = a.config.Client.Meta

	c.LogLevel = a.config.LogLevel
	c.Logger = a.logger

	return c, nil
}

func (a *Agent) setupHTTPServer() error {

	config := &http.Config{
		Logger:      a.logger,
		BindAddress: fmt.Sprintf("%s:%d", a.config.BindAddr, a.config.Ports.HTTP),
		Handlers: map[string]http.Handler{
			"/api/agent/":        handler.NewAgentHandler(a.rpcConn, a),
			"/api/nodes/":        handler.NewNodeHandler(a.rpcConn),
			"/api/interfaces/":   handler.NewInterfaceHandler(a.rpcConn),
			"/api/connections/":  handler.NewConnectionHandler(a.rpcConn),
			"/api/networks/":     handler.NewNetworkHandler(a.rpcConn),
			"/api/acl/":          handler.NewACLHandler(a.rpcConn),
			"/api/acl/tokens/":   handler.NewACLTokenHandler(a.rpcConn),
			"/api/acl/policies/": handler.NewACLPolicyHandler(a.rpcConn),
			"/status":            handler.NewStatusHandler(a.rpcConn),
		},
		Middleware: []http.Middleware{
			middleware.CORS(),
			middleware.Logging(a.logger),
		},
	}

	if a.config.UI {

		spaFS := a.config.StaticFS
		if !fsContains(spaFS, "index.html") {
			spaFS = nil
		}

		config.Handlers["/ui/"] = handler.NewSinglePageApplicationHandler(spaFS, uiStubHTML)
		config.Handlers["/"] = handler.NewFallthroughHandler("/ui/")
	}

	httpServer, err := http.NewServer(config)
	if err != nil {
		return err
	}

	a.httpServer = httpServer

	return nil
}

func fsContains(fs stdhttp.FileSystem, s string) bool {
	if _, err := fs.Open("index.html"); err != nil {
		return false
	}
	return true
}
