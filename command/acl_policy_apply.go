package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// ACLPolicyApplyCommand :
type ACLPolicyApplyCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	description string
}

func (c *ACLPolicyApplyCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	return flags
}

// Name :
func (c *ACLPolicyApplyCommand) Name() string {
	return "acl policy apply"
}

// Synopsis :
func (c *ACLPolicyApplyCommand) Synopsis() string {
	return "Apply ACL policy"
}

// Run :
func (c *ACLPolicyApplyCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <name>")
		c.UI.Error(`For additional help, try 'drago acl policy apply --help'`)
		return 1
	}

	name := args[0]

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	p := &structs.ACLPolicy{
		Name:        name,
		Description: c.description,
	}
	if err := api.ACLPolicies().Upsert(p); err != nil {
		c.UI.Error(fmt.Sprintf("Error applying ACL policy: %s", err))
		return 1
	}

	return 0
}

// Help :
func (c *ACLPolicyApplyCommand) Help() string {
	h := `
Usage: drago acl policy apply <name> [options]

  Create or update an ACL policy.

General Options:
` + GlobalOptions() + `

ACL Policy Apply Options:

  --description=<description>
    Sets the description for the ACL policy.

`
	return strings.TrimSpace(h)
}
