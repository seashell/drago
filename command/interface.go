package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// InterfaceCommand :
type InterfaceCommand struct {
	UI cli.UI
}

// Name :
func (c *InterfaceCommand) Name() string {
	return "interface"
}

// Synopsis :
func (c *InterfaceCommand) Synopsis() string {
	return "Interact with interfaces"
}

// Run :
func (c *InterfaceCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}

// Help :
func (c *InterfaceCommand) Help() string {
	h := `
Usage: drago interface <subcommand> [options] [args]

  This command groups subcommands for interacting with interfaces.
    
  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(h)
}
