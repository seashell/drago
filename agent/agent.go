package agent

import (
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"sync"

	banner "github.com/dimiro1/banner"
	application "github.com/seashell/drago/agent/application"
	http "github.com/seashell/drago/agent/infrastructure/http"
	handler "github.com/seashell/drago/agent/infrastructure/http/handler"
	client "github.com/seashell/drago/client"
	drago "github.com/seashell/drago/drago"
	log "github.com/seashell/drago/pkg/log"
	ui "github.com/seashell/drago/ui"
)

type Agent struct {
	config *Config

	logger log.Logger

	// Launched Drago Client, or nil if the agent isn't
	// configured to run a client.
	client *client.Client

	// Launched Drago Server, or nil if the agent isn't
	// configured to run a server.
	server *drago.Server

	// rpcClient
	rpcClient *rpc.Client

	// Application services exposed by the agent
	services struct {
		networks application.NetworkService
	}

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New
func New(config *Config, logger log.Logger) (*Agent, error) {

	config = DefaultConfig().Merge(config)

	a := &Agent{
		config:     config,
		logger:     logger,
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

	// Setup agent's RPC client/gateway
	if err := a.setupRPCClient(); err != nil {
		return nil, err
	}

	// Initialize agent application services
	a.services.networks = application.NewNetworkService(a.rpcClient)

	a.displayBanner()

	// Setup HTTP server
	if err := a.setupHTTPServer(); err != nil {
		return nil, err
	}

	return a, nil
}

// Shutdown
func (a *Agent) Shutdown() {
	a.logger.Debugf("Agent is shutting down...")
	// TODO: gracefully shutdown agent
}

// setupRPCClient
func (a *Agent) setupRPCClient() error {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:9999")
	if err != nil {
		return err
	}
	a.rpcClient = client
	return nil
}

// setupHTTPServer
func (a *Agent) setupHTTPServer() error {

	config := &http.ServerConfig{
		Logger:      a.logger,
		BindAddress: fmt.Sprintf("%s:%d", a.config.BindAddr, a.config.Ports.HTTP),
		Handlers: map[string]http.HandlerAdapter{
			"/api/healthcheck/": handler.NewHealthcheckHandlerAdapter(a.logger),
			"/api/networks/":    handler.NewNetworkHandlerAdapter(a.services.networks, a.logger),
		},
	}

	if a.config.UI {
		config.Handlers["/ui/"] = handler.NewFilesystemHandlerAdapter(ui.Bundle)
		config.Handlers["/"] = handler.NewFallthroughHandlerAdapter("/ui/")
	}

	httpServer, err := http.NewServer(config)
	if err != nil {
		return err
	}

	httpServer.Run()

	return nil
}

// displayBanner prints an ASCII banner to the standard output
func (a *Agent) displayBanner() {
	banner.Init(os.Stdout, true, true, strings.NewReader(bannerTmpl))
}

// setupServer
func (a *Agent) setupServer() error {

	if !a.config.Server.Enabled {
		return nil
	}

	config := &drago.Config{
		DataDir:  a.config.DataDir,
		BindAddr: fmt.Sprintf("%s:%d", a.config.BindAddr, a.config.Ports.RPC),
	}
	server, err := drago.NewServer(config)
	if err != nil {
		return fmt.Errorf("server setup failed: %v", err)
	}
	a.server = server

	return nil
}

// setupClient
func (a *Agent) setupClient() error {
	if !a.config.Client.Enabled {
		return nil
	}

	conf := &client.Config{}

	client, err := client.New(conf, a.logger)
	if err != nil {
		return fmt.Errorf("client setup failed: %v", err)
	}
	a.client = client

	return nil
}
