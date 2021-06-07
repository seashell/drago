package command

import (
	"context"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// NetworkDeleteCommand :
type NetworkDeleteCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json bool
}

func (c *NetworkDeleteCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *NetworkDeleteCommand) Name() string {
	return "network delete"
}

// Synopsis :
func (c *NetworkDeleteCommand) Synopsis() string {
	return "Delete an existing network"
}

// Run :
func (c *NetworkDeleteCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <network>")
		c.UI.Error(`For additional help, try 'drago network delete --help'`)
		return 1
	}

	name := args[0]
	id := ""

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	networks, err := api.Networks().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting networks: %s", err))
		return 1
	}

	for _, n := range networks {
		if n.Name == name {
			id = n.ID

			break
		}
	}

	if id == "" {
		c.UI.Error("Error: network not found")
		return 1
	}

	err = api.Networks().Delete(id)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting network: %s", err))
		return 1
	}

	return 0
}

// Help :
func (c *NetworkDeleteCommand) Help() string {
	h := `
Usage: drago network delete <network> [options]

  Delete an existing Drago network.

  If ACLs are enabled, this option requires a token with the 'network:write' capability.

General Options:
` + GlobalOptions()

	return strings.TrimSpace(h)
}
