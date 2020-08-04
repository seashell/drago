package command

import (
	"context"
	"strings"
)

type TokenListCommand struct {
	BaseCommand
}

func (c *TokenListCommand) Name() string {
	return "token list"
}

func (c *TokenListCommand) Synopsis() string {
	return "List all drago tokens"
}

func (c *TokenListCommand) Help() string {
	help := `
Usage: drago token list [options] [args]
		
  List all drago tokens

  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *TokenListCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	c.UI.Output("Not implemented\n")

	return 0
}
