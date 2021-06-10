package command

import (
	"context"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// NodeLeaveCommand :
type NodeLeaveCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	node string
}

func (c *NodeLeaveCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.node, "node", "", "")

	return flags
}

// Name :
func (c *NodeLeaveCommand) Name() string {
	return "node leave"
}

// Synopsis :
func (c *NodeLeaveCommand) Synopsis() string {
	return "Leave a network"
}

// Run :
func (c *NodeLeaveCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <network>")
		c.UI.Error(`For additional help, try 'drago node leave --help'`)
		return 1
	}

	networkID := args[0]
	nodeID := c.node

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	if nodeID == "" {
		if nodeID, err = localAgentNodeID(api); err != nil {
			c.UI.Error(fmt.Sprintf("Error determining local node ID: %s", err))
			return 1
		}
	}

	if _, err = api.Interfaces().Create(nodeID, networkID); err != nil {
		c.UI.Error(fmt.Sprintf("Error joining network: %s", err))
		return 1
	}

	c.UI.Output("Left!")

	return 0
}

// Help :
func (c *NodeLeaveCommand) Help() string {
	h := `
Usage: drago node join <network> [options]

  Have the local client node leave a network.

  If ACLs are enabled, this option requires a token with the 'interface:write' capability.

General Options:
` + GlobalOptions()

	return strings.TrimSpace(h)
}
