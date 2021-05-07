package command

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
)

// ACLPolicyUpsertCommand :
type ACLPolicyUpsertCommand struct {
	UI cli.UI

	// Parsed flags
	json bool

	Command
}

func (c *ACLPolicyUpsertCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ACLPolicyUpsertCommand) Name() string {
	return "acl policy upsert"
}

// Synopsis :
func (c *ACLPolicyUpsertCommand) Synopsis() string {
	return "Upsert ACL policy"
}

// Run :
func (c *ACLPolicyUpsertCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes no arguments")
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	p := &structs.ACLPolicy{}
	if err := api.ACLPolicies().Upsert(p); err != nil {
		c.UI.Error(fmt.Sprintf("Error upserting ACL policy: %s", err))
		return 1
	}

	c.UI.Output("Success!")

	return 0
}

// Help :
func (c *ACLPolicyUpsertCommand) Help() string {
	h := `
Usage: drago acl policy upsert [options]

  Create or update an ACL policy.

General Options:
` + GlobalOptions() + `

ACL Policy List Options:

  -json=<bool>
    Enable JSON output.

 `
	return strings.TrimSpace(h)
}
