package command

import (
	"context"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// ACLTokenDeleteCommand :
type ACLTokenDeleteCommand struct {
	UI cli.UI

	Command
}

func (c *ACLTokenDeleteCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	return flags
}

// Name :
func (c *ACLTokenDeleteCommand) Name() string {
	return "acl token delete"
}

// Synopsis :
func (c *ACLTokenDeleteCommand) Synopsis() string {
	return "Delete an existing ACL token"
}

// Run :
func (c *ACLTokenDeleteCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.UI.Error("This command takes one argument: <id>")
		c.UI.Error(`For additional help, try 'drago acl token delete --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	if err := api.ACLTokens().Delete(args[0]); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting ACL token: %s", err))
		return 1
	}

	c.UI.Output("ACL token successfully deleted")

	return 0
}

// Help :
func (c *ACLTokenDeleteCommand) Help() string {
	h := `
Usage: drago acl token delete <token> [options]

  Delete is used to delete an existing ACL token. Requires a management token.

General Options:
` + GlobalOptions()

	return strings.TrimSpace(h)
}
