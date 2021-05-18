package command

import (
	"context"
	"flag"
	"fmt"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
)

// InterfaceUpdateCommand :
type InterfaceUpdateCommand struct {
	UI cli.UI

	// Parsed flags
	address string

	Command
}

func (c *InterfaceUpdateCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.address, "address", "", "")

	return flags
}

// Name :
func (c *InterfaceUpdateCommand) Name() string {
	return "interface update"
}

// Synopsis :
func (c *InterfaceUpdateCommand) Synopsis() string {
	return "Update an existing interfaces"
}

// Run :
func (c *InterfaceUpdateCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 0 {
		c.UI.Error("This command takes one argument")
		return 1
	}

	id := args[0]

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	api.Interfaces().Update(&structs.Interface{
		ID:      id,
		Address: &c.address,
	})

	return 0
}

// Help :
func (c *InterfaceUpdateCommand) Help() string {
	h := `
Usage: drago interface update <id> [options]

  Update an existing interface.

  If ACLs are enabled, this option requires a token with the 'interface:write' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  -address=<id>
    Interface address.

 `
	return strings.TrimSpace(h)
}
