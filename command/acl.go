package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// ACLCommand :
type ACLCommand struct {
	UI cli.UI
}

// Name :
func (c *ACLCommand) Name() string {
	return "acl"
}

// Synopsis :
func (c *ACLCommand) Synopsis() string {
	return "Interact with ACL policies and tokens"
}

// Run :
func (c *ACLCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *ACLCommand) Help() string {
	h := `
Usage: drago acl <subcommand> [options] [args]

  This command groups subcommands for interacting with ACL policies and tokens.
  It can be used to bootstrap Drago's ACL system, create policies that restrict
  access, and generate tokens from those policies.
  
  Bootstrap ACLs:

      $ drago acl bootstrap

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(h)
}
