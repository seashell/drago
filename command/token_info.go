package command

import (
	"context"
	"strings"
)

type TokenInfoCommand struct {
	BaseCommand
}

func (c *TokenInfoCommand) Name() string {
	return "token info"
}

func (c *TokenInfoCommand) Synopsis() string {
	return "Get information about a drago token"
}

func (c *TokenInfoCommand) Help() string {
	help := `
Usage: drago token info [options] [args]
		
  Get information about a drago token

  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *TokenInfoCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	c.UI.Output("Not implemented\n")

	return 0
}
