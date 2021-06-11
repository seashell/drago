package command

import (
	"context"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// ConnectionDeleteCommand :
type ConnectionDeleteCommand struct {
	UI cli.UI
	Command
}

func (c *ConnectionDeleteCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	return flags
}

// Name :
func (c *ConnectionDeleteCommand) Name() string {
	return "connection delete"
}

// Synopsis :
func (c *ConnectionDeleteCommand) Synopsis() string {
	return "Delete an existing connection"
}

// Run :
func (c *ConnectionDeleteCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 3 {
		c.UI.Error("This command takes one argument: <connection_id>")
		c.UI.Error(`For additional help, try 'drago connection delete --help'`)
		return 1
	}

	id := args[0]

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	err = api.Connections().Delete(id)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting connection: %s", err))
		return 1
	}

	c.UI.Output("Deleted!")

	return 0
}

// Help :
func (c *ConnectionDeleteCommand) Help() string {
	h := `
Usage: drago connection delete <connection_id> [options]

  Delete is used to delete an existing connection.

  If ACLs are enabled, this option requires a token with the 'connection:write' capability.

General Options:
` + GlobalOptions() + `
`
	return strings.TrimSpace(h)
}
