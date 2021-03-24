package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// NodeCommand :
type NodeCommand struct {
	UI cli.UI
}

// Name :
func (c *NodeCommand) Name() string {
	return "node"
}

// Synopsis :
func (c *NodeCommand) Synopsis() string {
	return "Interact with nodes"
}

// Run :
func (c *NodeCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *NodeCommand) Help() string {
	h := `
Usage: drago node <subcommand> [options] [args]

  This command groups subcommands for interacting with nodes.
    
  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(h)
}
