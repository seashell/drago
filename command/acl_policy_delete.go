package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	cli "github.com/seashell/drago/pkg/cli"
)

// ACLPolicyDeleteCommand :
type ACLPolicyDeleteCommand struct {
	UI cli.UI
	Command
}

func (c *ACLPolicyDeleteCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	return flags
}

// Name :
func (c *ACLPolicyDeleteCommand) Name() string {
	return "acl policy delete"
}

// Synopsis :
func (c *ACLPolicyDeleteCommand) Synopsis() string {
	return "Delete ACL policy"
}

// Run :
func (c *ACLPolicyDeleteCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <name>")
		c.UI.Error(`For additional help, try 'drago acl policy delete --help'`)
		return 1
	}

	name := args[0]

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	if err := api.ACLPolicies().Delete(name); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting ACL policies: %s", err))
		return 1
	}

	return 0
}

// Help :
func (c *ACLPolicyDeleteCommand) Help() string {
	h := `
Usage: drago acl policy delete <name> [options]

  Delete an existing ACL policy.

General Options:
` + GlobalOptions()

	return strings.TrimSpace(h)
}
