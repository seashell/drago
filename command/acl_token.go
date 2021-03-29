package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// ACLTokenCommand :
type ACLTokenCommand struct {
	UI cli.UI
}

// Name :
func (c *ACLTokenCommand) Name() string {
	return "acl token"
}

// Synopsis :
func (c *ACLTokenCommand) Synopsis() string {
	return "Interact with ACL tokens"
}

// Run :
func (c *ACLTokenCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *ACLTokenCommand) Help() string {
	h := `
Usage: drago acl token <subcommand> [options] [args]

  This command groups subcommands for interacting with ACL tokens.
  
  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(h)
}
