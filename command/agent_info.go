package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// AgentInfoCommand :
type AgentInfoCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json bool
}

func (c *AgentInfoCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *AgentInfoCommand) Name() string {
	return "agent-info"
}

// Synopsis :
func (c *AgentInfoCommand) Synopsis() string {
	return "Display status information about the local agent"
}

// Run :
func (c *AgentInfoCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes no arguments")
		c.UI.Error(`For additional help, try 'drago agent info --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	info, err := api.Agent().Self()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving agent info: %s", err))
		return 1
	}

	c.UI.Output(c.formatAgentInfo(info))

	return 0
}

// Help :
func (c *AgentInfoCommand) Help() string {
	h := `
Usage: drago agent-info [options]

  Display status information about the local agent.

General Options:
` + GlobalOptions() + `

Agent Info Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *AgentInfoCommand) formatAgentInfo(info *structs.Agent) string {

	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")

	fpolicy := map[string]interface{}{
		"config": info.Config,
		"stats":  info.Stats,
	}
	if err := enc.Encode(fpolicy); err != nil {
		c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
	}

	s := b.String()

	if c.json {
		return s
	}

	return cleanJSONString(s)
}
