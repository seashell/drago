package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

type HostCommand struct {
	BaseCommand
}

func (c *HostCommand) Name() string {
	return "host"
}

func (c *HostCommand) Synopsis() string {
	return "Interacts with drago hosts"
}

func (c *HostCommand) Help() string {
	txt := `
Usage: drago host <subcomand> [options] [args]
		
  This command groups subcommands for interacting with hosts. 

  Please see the individual subcommand help for detailed usage information.

`
	return strings.TrimSpace(txt)
}

func (c *HostCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}
