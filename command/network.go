package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// NetworkCommand :
type NetworkCommand struct {
	UI cli.UI
}

// Name :
func (c *NetworkCommand) Name() string {
	return "network"
}

// Synopsis :
func (c *NetworkCommand) Synopsis() string {
	return "Interact with networks"
}

// Run :
func (c *NetworkCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *NetworkCommand) Help() string {
	h := `
Usage: drago network <subcommand> [options] [args]

  This command groups subcommands for interacting with networks.
    
  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(h)
}
