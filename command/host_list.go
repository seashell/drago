package command

import (
	"context"
	"strings"
)

type HostListCommand struct {
	BaseCommand
}

func (c *HostListCommand) Name() string {
	return "host list"
}

func (c *HostListCommand) Synopsis() string {
	return "List all drago hosts"
}

func (c *HostListCommand) Help() string {
	help := `
Usage: drago host list [options] [args]
		
  List hosts.
  
  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *HostListCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	c.UI.Output("Not implemented\n")

	return 0
}
