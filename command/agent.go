package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	agent "github.com/seashell/drago/agent"
	log "github.com/seashell/drago/pkg/log"
	logrus "github.com/seashell/drago/pkg/log/logrus"
	version "github.com/seashell/drago/version"
)

type AgentCommand struct {
	BaseCommand
}

func (c *AgentCommand) Name() string {
	return "agent"
}

func (c *AgentCommand) Synopsis() string {
	return "Runs a drago agent"
}

func (c *AgentCommand) Run(ctx context.Context, args []string) int {

	config := c.parseConfig(args)

	logger, err := logrus.NewLoggerAdapter(logrus.Config{
		LoggerOptions: log.LoggerOptions{
			Prefix: "agent: ",
			Level:  logrus.Debug,
		},
	})
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	agent, err := agent.New(config, logger)
	if err != nil {
		c.UI.Error("==> " + "Error initializing agent: " + err.Error() + "\n")
		return 1
	}

	defer func() {
		agent.Shutdown()
	}()

	c.printAgentConfig(config)

	c.UI.Output("==> Drago agent started! Log data will stream in below:\n")

	return c.handleSignals()
}

// parseConfig
func (c *AgentCommand) parseConfig(args []string) *agent.Config {

	var devMode bool
	var configPath string

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	cmdConfig := &agent.Config{
		Server:  &agent.ServerConfig{},
		Client:  &agent.ClientConfig{},
		ACL:     &agent.ACLConfig{},
		Ports:   &agent.Ports{},
		Version: version.GetVersion(),
	}

	// Agent mode
	flags.BoolVar(&devMode, "dev", false, "")
	flags.BoolVar(&cmdConfig.Server.Enabled, "server", false, "")
	flags.BoolVar(&cmdConfig.Client.Enabled, "client", false, "")

	// Client-only options
	flags.StringVar(&cmdConfig.Client.StateDir, "state-dir", "", "")
	flags.StringVar(&cmdConfig.Client.InterfacesPrefix, "interface-prefix", "", "")
	flags.StringVar(&cmdConfig.Client.Server, "remote-server", "", "")

	// Server-only options
	// --

	// General options
	flags.StringVar(&configPath, "config", "", "")
	flags.StringVar(&cmdConfig.BindAddr, "bind", "", "")
	flags.StringVar(&cmdConfig.DataDir, "data_dir", "", "")
	flags.StringVar(&cmdConfig.LogLevel, "log-level", "", "")
	flags.StringVar(&cmdConfig.PluginDir, "plugin-dir", "", "")

	// ACL options
	flags.BoolVar(&cmdConfig.ACL.Enabled, "acl-enabled", false, "")

	if err := flags.Parse(args); err != nil {
		c.UI.Error("==> Error: " + err.Error() + "\n")
		return nil
	}

	config := agent.DefaultConfig()

	if configPath != "" {
		config := &agent.Config{}
		config, err := config.LoadFromFile(configPath)
		if err != nil {
			c.UI.Error(fmt.Sprintf("==> Error parsing configuration from file: %s", err))
		}
		c.UI.Output(fmt.Sprintf("==> Loaded configuration from %s", configPath))
	} else {
		c.UI.Warn("==> No configuration file loaded")
	}

	config = config.Merge(cmdConfig)
	if !config.IsValid() {
		return nil
	}

	return config
}

// printAgentConfig
func (c *AgentCommand) printAgentConfig(config *agent.Config) {

	info := map[string]string{
		"bind addr": config.BindAddr,
		"client":    strconv.FormatBool(config.Client.Enabled),
		"log level": config.LogLevel,
		"server":    strconv.FormatBool(config.Server.Enabled),
		"version":   config.Version.VersionNumber(),
	}

	padding := 18
	c.UI.Output("==> Drago agent configuration:\n")
	for k, _ := range info {
		c.UI.Info(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")
}

// handleSignals waits for specific signals and returns
func (c *AgentCommand) handleSignals() int {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Wait until a signal is received
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	}

	c.UI.Output(fmt.Sprintf("Caught signal: %v", sig))

	return 1
}

func (c *AgentCommand) Help() string {
	h := `
Usage: drago agent [options]
	
  Starts the Drago agent and runs until an interrupt is received.
  The agent may be a client and/or server.
  
  The Drago agent's configuration primarily comes from the config
  files used, but a subset of the options may also be passed directly
  as CLI arguments.

General Options:

  -bind=<addr>
    The address the agent will bind to for all of its various network
    services. The individual services that run bind to individual
    ports on this address. Defaults to 127.0.0.1.

  -data-dir=<path>
    The data directory where all state will be persisted. On Drago 
    clients this is used to store local network configurations, whereas
    on server nodes, the data dir is also used to keep the desired state
    for all the managed networks.

  -plugin-dir=<path>
    The plugin directory from which Drago plugins will be loaded.
    If not specified, the plugin directory defaults <data-dir>/plugins/.

  -config=<path>
    The path to a config file to use for configuring the Drago agent.

  -dev=<bool>
    Start the agent in development mode, which means that both the
    Drago client and server will be active at the same time. This
    is useful for development and testing purposes. If the --dev flag
    is passed to the agent, no other configuration is required.
	
  -log-level=<level>
    Specify the verbosity level of Drago's logs. Valid values include
    DEBUG, INFO, WARN, ERROR, and FATAL in decreasing order of verbosity.
    The	default is INFO.

Server Options:

  -server=<bool>
    Start the agent in server mode. A Drago server is responsible for
    interacting with the storage backend and providing clients with
    up-to-date configurations for their network interfaces.
  
Client Options:
  -client=<bool>
    Start the agent in client mode. A Drago client is responsible for
    fetching the desired networking state from the server, and performing
    a reconciliation procedure on its network interfaces. Additionally,
    it is responsible for updating the seerver of its own state, including
    metrics, and currently active public keys.

  -state-dir
    The directory used to store state and other persistent data. If not
    specified a subdirectory under the higher-level option "-data-dir"
    will be used.
	
  -interface-prefix=<string>
    Specify the prefix used to identify all network interfaces managed by
    Drago on client machines. This allow Drago to selectively apply configurations
    and alter the state of specific interfaces without affecting others
    that might have been created and are managed by the user.

  -remote-server=<string>
    The address of a known Drago server in format "host:port" with which the client
    will connect to obtain its network configuration.

ACL Options:

  -acl-enabled
    Specifies whether the agent should enable ACLs.


`
	return strings.TrimSpace(h)
}
