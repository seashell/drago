package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

type TokenCommand struct {
	BaseCommand
}

func (c *TokenCommand) Name() string {
	return "token"
}

func (c *TokenCommand) Synopsis() string {
	return "Interacts with drago tokens"
}

func (c *TokenCommand) Help() string {
	txt := `
Usage: drago token <subcomand> [options] [args]
		
  This command groups subcommands for interacting with tokens. 

  Please see the individual subcommand help for detailed usage information.

`
	return strings.TrimSpace(txt)
}

func (c *TokenCommand) Run(ctx context.Context, args []string) int {
	return cli.CommandReturnCodeHelp
}
