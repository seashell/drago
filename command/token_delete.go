package command

import (
	"context"
	"strings"
)

type TokenDeleteCommand struct {
	BaseCommand
}

func (c *TokenDeleteCommand) Name() string {
	return "token delete"
}

func (c *TokenDeleteCommand) Synopsis() string {
	return "Delete a drago token"
}

func (c *TokenDeleteCommand) Help() string {
	help := `
Usage: drago token delete [options] [args]
		
  Delete a drago token.

  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *TokenDeleteCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	c.UI.Output("Not implemented\n")

	return 0
}
