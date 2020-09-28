package command

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
	"github.com/dimiro1/banner"
	"github.com/joho/godotenv"
	agent "github.com/seashell/drago/agent"
	cli "github.com/seashell/drago/pkg/cli"
	log "github.com/seashell/drago/pkg/log"
	zap "github.com/seashell/drago/pkg/log/zap"
)

// AgentCommand :
type AgentCommand struct {
	UI cli.UI
}

// Name :
func (c *AgentCommand) Name() string {
	return "agent"
}

// Synopsis :
func (c *AgentCommand) Synopsis() string {
	return "Runs a drago agent"
}

// Run :
func (c *AgentCommand) Run(ctx context.Context, args []string) int {

	displayBanner()

	config := c.parseConfig(args)

	// logger, err := logrus.NewLoggerAdapter(logrus.Config{
	// 	LoggerOptions: log.LoggerOptions{
	// 		Level:  config.LogLevel,
	// 		Prefix: "agent: ",
	// 	},
	// })

	logger, err := zap.NewLoggerAdapter(zap.Config{
		LoggerOptions: log.LoggerOptions{
			Level:  config.LogLevel,
			Prefix: "agent: ",
		},
	})

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output("==> Starting Drago agent...")

	// Create DataDir and other subdirectories if they do not exist
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		os.Mkdir(config.DataDir, 0755)
	}

	c.printConfig(config)

	_, err = agent.New(config, logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error starting agent: %s\n", err.Error()))
	}

	<-ctx.Done()

	return 0
}

func (c *AgentCommand) parseConfig(args []string) *agent.Config {

	flags := FlagSet(c.Name())

	configFromFlags := c.parseFlags(flags, args)
	configFromFile := c.parseConfigFiles(flags.configPaths...)
	configFromEnv := c.parseEnv(flags.envPaths...)

	config := agent.DefaultConfig()

	config = config.Merge(configFromFile)
	config = config.Merge(configFromEnv)
	config = config.Merge(configFromFlags)

	if err := config.Validate(); err != nil {
		c.UI.Error(fmt.Sprintf("Invalid input: %s", err.Error()))
		os.Exit(1)
	}

	return config
}

func (c *AgentCommand) parseFlags(flags *RootFlagSet, args []string) *agent.Config {

	flags.Usage = func() {
		c.UI.Output("\n" + c.Help() + "\n")
	}

	config := agent.EmptyConfig()

	var devMode bool

	// Agent mode
	flags.BoolVar(&devMode, "dev", false, "")
	flags.BoolVar(&config.Server.Enabled, "server", false, "")
	flags.BoolVar(&config.Client.Enabled, "client", false, "")

	// General options (available in both client and server modes)
	flags.StringVar(&config.DataDir, "data-dir", "", "")
	flags.StringVar(&config.BindAddr, "bind-addr", "", "")
	flags.StringVar(&config.PluginDir, "plugin-dir", "", "")
	flags.StringVar(&config.LogLevel, "log-level", "", "")

	// Client-only options
	flags.StringVar(&config.Client.StateDir, "state-dir", "", "")

	// Server-only options
	// --

	// ACL options
	flags.BoolVar(&config.ACL.Enabled, "acl-enabled", false, "")

	if err := flags.Parse(args); err != nil {
		c.UI.Error("==> Error: " + err.Error() + "\n")
		os.Exit(1)
	}

	return config
}

func (c *AgentCommand) parseConfigFiles(paths ...string) *agent.Config {

	config := agent.EmptyConfig()

	if len(paths) > 0 {
		// TODO : Load configurations from HCL files
		c.UI.Info(fmt.Sprintf("==> Loading configurations from: %v", paths))
	} else {
		c.UI.Output("==> No configuration files loaded")
	}

	return config
}

func (c *AgentCommand) parseEnv(paths ...string) *agent.Config {

	config := agent.EmptyConfig()

	if len(paths) > 0 {

		c.UI.Info(fmt.Sprintf("==> Loading environment variables from: %v", paths))
		c.UI.Warn(fmt.Sprintf("  - This will not override already existing variables!"))

		err := godotenv.Load(paths...)

		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing env files: %s", err.Error()))
			os.Exit(1)
		}
	}

	env.Parse(config)

	return config
}

func (c *AgentCommand) printConfig(config *agent.Config) {

	info := map[string]string{
		"data dir":        config.DataDir,
		"bind addrs":      bindAddrsString(config),
		"advertise addrs": advertiseAddrsString(config),
		"log level":       config.LogLevel,
		"client":          strconv.FormatBool(config.Client.Enabled),
		"server":          strconv.FormatBool(config.Server.Enabled),
		"version":         config.Version.VersionNumber(),
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

// Help :
func (c *AgentCommand) Help() string {
	h := `
Usage: drago agent [options]
	
  Starts the Drago agent and runs until an interrupt is received.
  The agent may be a client and/or server.
  
  The Drago agent's configuration primarily comes from the config
  files used, but a subset of the options may also be passed directly
  as CLI arguments.

General Options:
` + GlobalOptions() + `

Agent Options:

  --data-dir=<path>
    The data directory where all state will be persisted. On Drago 
    clients this is used to store local network configurations, whereas
    on server nodes, the data dir is also used to keep the desired state
	for all the managed networks. Overrides the DRAGO_DATA_DIRenvironment
	variable if set.

  --log-level=<level>
    The logging level Drago should log at. Valid values are INFO, WARN, DEBUG, ERROR, FATAL.
    Overrides the DRAGO_LOG_LEVEL environment variable if set.
	
`
	return strings.TrimSpace(h)
}

// Prints an ASCII banner to the standard output
func displayBanner() {
	banner.Init(os.Stdout, true, true, strings.NewReader(agent.Banner))
}

func bindAddrsString(config *agent.Config) string {
	http := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.HTTP)
	rpc := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.RPC)
	return fmt.Sprintf("HTTP: %s; RPC: %s", http, rpc)
}

func advertiseAddrsString(config *agent.Config) string {

	http := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.HTTP)
	//if config.AdvertiseAddrs.Peer != "" {
	//	http = fmt.Sprintf("%s", config.AdvertiseAddrs.HTTP)
	//}

	rpc := fmt.Sprintf("%s:%d", config.BindAddr, config.Ports.RPC)
	//if config.AdvertiseAddrs.HTTP != "" {
	//	rpc = fmt.Sprintf("%s", config.AdvertiseAddrs.RPC)
	//}

	return fmt.Sprintf("HTTP: %s; RPC: %s", http, rpc)
}
