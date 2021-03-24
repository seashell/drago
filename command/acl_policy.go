package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// ACLPolicyCommand :
type ACLPolicyCommand struct {
	UI cli.UI
}

// Name :
func (c *ACLPolicyCommand) Name() string {
	return "acl policy"
}

// Synopsis :
func (c *ACLPolicyCommand) Synopsis() string {
	return "Interact with ACL policies"
}

// Run :
func (c *ACLPolicyCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *ACLPolicyCommand) Help() string {
	h := `
Usage: drago acl policy <subcommand> [options] [args]

  This command groups subcommands for interacting with ACL policies.

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(h)
}
