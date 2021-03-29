package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// ConnectionCommand :
type ConnectionCommand struct {
	UI cli.UI
}

// Name :
func (c *ConnectionCommand) Name() string {
	return "connection"
}

// Synopsis :
func (c *ConnectionCommand) Synopsis() string {
	return "Interact with connections"
}

// Run :
func (c *ConnectionCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *ConnectionCommand) Help() string {
	h := `
Usage: drago connection <subcommand> [options] [args]

  This command groups subcommands for interacting with connections.
    
  Please see the individual subcommand help for detailed usage information.

General Options:
` + GlobalOptions() + `
`
	return strings.TrimSpace(h)
}
