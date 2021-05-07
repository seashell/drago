package command

import (
	"context"
	"flag"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// NetworkDeleteCommand :
type NetworkDeleteCommand struct {
	UI cli.UI

	// Parsed flags
	json bool
	name string

	Command
}

func (c *NetworkDeleteCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")
	flags.StringVar(&c.name, "name", "", "")

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
	if len(args) > 1 {
		c.UI.Error("This command takes either one or no arguments")
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	id := ""

	if len(args) > 0 {
		id = args[0]
	}

	if c.name != "" {
		networks, err := api.Networks().List()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting networks: %s", err))
			return 1
		}

		for _, n := range networks {
			if n.Name == c.name {
				if id != "" && n.ID != id {
					c.UI.Error("Error: name and ID belong to different networks")
					return 1
				}
				id = n.ID
			}
		}
	}

	if id == "" {
		c.UI.Error("Error: network not found")
		return 1
	}

	err = api.Networks().Delete(id)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating network: %s", err))
		return 1
	}

	return 0
}

// Help :
func (c *NetworkDeleteCommand) Help() string {
	h := `
Usage: drago network delete <id> [options]

  Delete an existing Drago network.

  If ACLs are enabled, this option requires a token with the 'network:write' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  -name=""
    Human readable name of the network to be deleted.
 `
	return strings.TrimSpace(h)
}
