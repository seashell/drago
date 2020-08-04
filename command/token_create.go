package command

import (
	"context"
	"strings"
)

type TokenCreateCommand struct {
	BaseCommand
}

func (c *TokenCreateCommand) Name() string {
	return "token create"
}

func (c *TokenCreateCommand) Synopsis() string {
	return "Create a drago token"
}

func (c *TokenCreateCommand) Help() string {
	help := `
Usage: drago token create [options] [args]
		
  Create a drago token. 

  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *TokenCreateCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	c.UI.Output("Not implemented\n")

	return 0
}
