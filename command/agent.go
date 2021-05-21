package command

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
	"github.com/dimiro1/banner"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/joho/godotenv"
	agent "github.com/seashell/drago/agent"
	cli "github.com/seashell/drago/pkg/cli"
	log "github.com/seashell/drago/pkg/log"
	simple "github.com/seashell/drago/pkg/log/simple"
	"github.com/spf13/pflag"
)

// AgentCommand :
type AgentCommand struct {
	config *agent.Config
	agent  *agent.Agent
	logger log.Logger

	UI       cli.UI
	StaticFS http.FileSystem

	Command

	// Parsed flags
	dev           bool
	servers       string
	envs          []string
	configs       []string
	meta          []string
	node          string
	bind          string
	dataDir       string
	logLevel      string
	pluginDir     string
	server        bool
	client        bool
	stateDir      string
	wireguardPath string
	aclEnabled    bool
}

func (c *AgentCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options (available in both client and server modes)
	flags.StringSliceVar(&c.configs, "config", []string{}, "")
	flags.StringVar(&c.bind, "bind", "", "")
	flags.StringVar(&c.node, "node", "", "")
	flags.BoolVar(&c.dev, "dev", false, "")
	flags.StringVar(&c.dataDir, "data-dir", "", "")
	flags.StringVar(&c.logLevel, "log-level", "", "")
	flags.StringVar(&c.pluginDir, "plugin-dir", "", "")
	flags.StringSliceVar(&c.envs, "env", []string{}, "")

	// Server-only options
	flags.BoolVar(&c.server, "server", false, "")

	// Client-only options
	flags.StringSliceVar(&c.meta, "meta", []string{}, "")
	flags.BoolVar(&c.client, "client", false, "")
	flags.StringVar(&c.servers, "servers", "", "")
	flags.StringVar(&c.stateDir, "state-dir", "", "")
	flags.StringVar(&c.wireguardPath, "wireguard-path", "", "")

	// ACL options
	flags.BoolVar(&c.aclEnabled, "acl-enabled", false, "")

	return flags
}

// Name :
func (c *AgentCommand) Name() string {
	return "agent"
}

// Synopsis :
func (c *AgentCommand) Synopsis() string {
	return "Run a Drago agent"
}

// Run :
func (c *AgentCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		c.UI.Error("==> Error: " + err.Error() + "\n")
		os.Exit(1)
	}

	if !c.server && !c.client && !c.dev {
		c.UI.Output("==> Must specify either client, server or dev mode for the agent.")
		os.Exit(1)
	}

	printBanner()

	err := c.parseConfig(args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Invalid input: %s", err.Error()))
		os.Exit(1)
	}

	if err := c.setupLogger(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output("==> Starting Drago agent...")

	if err = c.setupDirectories(); err != nil {
		c.UI.Error("Error setting up data directories")
	}

	c.printConfig()

	if err = c.setupAgent(); err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up agent: %s\n", err.Error()))
		return 1
	}

	<-ctx.Done()

	c.agent.Shutdown()

	return 0
}

// Help :
func (c *AgentCommand) Help() string {
	h := `
Usage: drago agent [options]
	
  Starts the Drago agent and runs until an interrupt is received.
  The agent may be a client and/or server.
  
  The Drago agent's configuration primarily comes from the config
  files used, but a subset of the options may also be passed directly
  as CLI arguments.

General Options (clients and servers):
` + GlobalOptions() + `

  --bind=<addr>
    The address the agent will bind to for all of its various network
    services. The individual services that run bind to individual
    ports on this address. Defaults to the loopback 127.0.0.1.

  --config=<path>
    Path to a HCL file containing valid Drago configurations.
    Overrides the DRAGO_CONFIG_PATH environment variable if set.

  --data-dir=<path>
    The data directory where all state will be persisted. On Drago 
    clients this is used to store local network configurations, whereas
    on server nodes, the data dir is also used to keep the desired state
    for all the managed networks. Overrides the DRAGO_DATA_DIRenvironment
    variable if set.

  --node=<name>
    The name of the local agent, use to identify the node. If not provided,
    defaults to the hostname of the machine.

  --dev
    Start the agent in development mode. This enables a pre-configured
    dual-role agent (client + server) which is useful for developing
    or testing Drago. No other configuration is required to start the
    agent in this mode.

  --log-level=<level>
    The logging level Drago should log at. Valid values are INFO, WARN,
    DEBUG, ERROR, FATAL. Overrides the DRAGO_LOG_LEVEL environment variable
    if set.

  --plugin-dir=<path>
    The plugin directory is used to discover Drago plugins. If not specified,
    the plugin directory defaults to be that of <data-dir>/plugins/.

Server Options:

  --server
    Enable server mode for the agent.

Client Options:

  --client
    Enable client mode for the agent. Client mode enables a given node to be
    evaluated for allocations. If client mode is not enabled, no work will be
    scheduled to the agent.

  --state-dir
    The directory used to store state and other persistent data. If not
    specified a subdirectory under the "-data-dir" will be used.
  
  --servers
    A comma-separated list of known server addresses to connect to in
    "host:port" format.

  --meta
    User-specified metadata in KEY=VALUE format to associate with the node.
	Repeat the meta flag for each key/value pair to be added.

ACL Options:

  --acl-enabled
    Specifies whether the agent should enable ACLs.

`
	return strings.TrimSpace(h)
}

func (c *AgentCommand) setupAgent() error {

	c.config.StaticFS = c.StaticFS

	agent, err := agent.New(c.config, c.logger)
	if err != nil {
		return err
	}

	c.agent = agent

	return nil
}

func (c *AgentCommand) parseConfig(args []string) error {

	configFromFlags := c.parseFlags()

	configFromFile := c.parseConfigFiles(c.configs...)
	configFromEnv := c.parseEnv(c.envs...)

	config := agent.DefaultConfig()

	config = config.Merge(configFromFile)
	config = config.Merge(configFromEnv)
	config = config.Merge(configFromFlags)

	if err := config.Validate(); err != nil {
		return err
	}

	c.config = config

	return nil
}

func (c *AgentCommand) parseFlags() *agent.Config {

	config := agent.EmptyConfig()

	config.DevMode = &c.dev
	config.Server.Enabled = c.server
	config.Client.Enabled = c.client

	config.Name = c.node
	config.BindAddr = c.bind
	config.DataDir = c.dataDir
	config.LogLevel = c.logLevel
	config.PluginDir = c.pluginDir

	config.Client.StateDir = c.stateDir
	config.Client.WireguardPath = c.wireguardPath

	config.ACL.Enabled = c.aclEnabled

	if *config.DevMode {
		config.Server.Enabled = true
		config.Client.Enabled = true
		config.DataDir = "/tmp/drago"
		config.LogLevel = "DEBUG"
	}

	if c.servers != "" {
		config.Client.Servers = strings.Split(c.servers, ",")
	}

	metaLength := len(c.meta)
	if metaLength != 0 {
		config.Client.Meta = make(map[string]string, metaLength)
		for _, kv := range c.meta {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				c.UI.Error(fmt.Sprintf("Error parsing Client.Meta value: %v", kv))
				return nil
			}
			config.Client.Meta[parts[0]] = parts[1]
		}
	}

	return config
}

func (c *AgentCommand) parseConfigFiles(paths ...string) *agent.Config {

	config := agent.EmptyConfig()

	if len(paths) > 0 {
		c.UI.Info(fmt.Sprintf("==> Loading configurations from: %v", paths))
		for _, s := range paths {
			err := hclsimple.DecodeFile(s, nil, config)
			if err != nil {
				c.UI.Error("Failed to load configuration: " + err.Error())
				os.Exit(0)
			}
		}
	} else {
		c.UI.Output("==> No configuration files loaded")
	}

	return config
}

func (c *AgentCommand) parseEnv(paths ...string) *agent.Config {

	config := agent.EmptyConfig()

	if len(paths) > 0 {

		c.UI.Info(fmt.Sprintf("==> Loading environment variables from: %v", paths))
		c.UI.Warn("  - This will not override already existing variables!")

		err := godotenv.Load(paths...)

		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing env files: %s", err.Error()))
			os.Exit(1)
		}
	}

	env.Parse(config)

	return config
}

func (c *AgentCommand) printConfig() {

	config := c.config

	info := map[string]string{
		"data dir":        config.DataDir,
		"bind addrs":      bindAddrsString(config),
		"advertise addrs": advertiseAddrsString(config),
		"log level":       config.LogLevel,
		"client":          strconv.FormatBool(config.Client.Enabled),
		"server":          strconv.FormatBool(config.Server.Enabled),
		"version":         config.Version.VersionNumber(),
		"acl enabled":     strconv.FormatBool(config.ACL.Enabled),
	}

	padding := 18
	c.UI.Output("==> Drago agent configuration:\n")
	for k := range info {
		c.UI.Info(fmt.Sprintf(
			"%s%s: %v",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}

	c.UI.Output("")
}

func (c *AgentCommand) setupLogger() error {

	// logger, err := logrus.NewLoggerAdapter(logrus.Config{
	// 	LoggerOptions: log.LoggerOptions{
	// 		Level:  c.config.LogLevel,
	// 		Prefix: "agent: ",
	// 	},
	// })

	// logger, err := zap.NewLoggerAdapter(zap.Config{
	// 	LoggerOptions: log.LoggerOptions{
	// 		Level:  c.config.LogLevel,
	// 		Prefix: "agent: ",
	// 	},
	// })

	logger, err := simple.NewLoggerAdapter(simple.Config{
		LoggerOptions: log.LoggerOptions{
			Level:  c.config.LogLevel,
			Prefix: "agent: ",
		},
	})

	if err != nil {
		return err
	}

	c.logger = logger

	return nil
}

// Create DataDir and other subdirectories if they do not exist
func (c *AgentCommand) setupDirectories() error {
	if _, err := os.Stat(c.config.DataDir); os.IsNotExist(err) {
		os.Mkdir(c.config.DataDir, 0755)
	}
	return nil
}

// Prints an ASCII banner to the standard output
func printBanner() {
	banner.Init(os.Stdout, true, true, strings.NewReader(agent.Banner))
}

func bindAddrsString(config *agent.Config) string {
	http := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.HTTP)
	rpc := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.RPC)
	return fmt.Sprintf("HTTP: %s; RPC: %s", http, rpc)
}

func advertiseAddrsString(config *agent.Config) string {

	peer := config.BindAddr
	if config.AdvertiseAddrs.Peer != "" {
		peer = config.AdvertiseAddrs.Peer
	}

	server := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.RPC)
	if config.AdvertiseAddrs.Server != "" {
		server = fmt.Sprintf("%s:%d", config.AdvertiseAddrs.Server, config.Ports.RPC)
	}

	return fmt.Sprintf("Peer: %s; Server: %s", peer, server)
}
