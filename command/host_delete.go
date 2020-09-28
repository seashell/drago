package command

import (
	"context"
	"strings"
)

type HostDeleteCommand struct {
	BaseCommand
}

func (c *HostDeleteCommand) Name() string {
	return "host delete"
}

func (c *HostDeleteCommand) Synopsis() string {
	return "Delete a drago host"
}

func (c *HostDeleteCommand) Help() string {
	help := `
Usage: drago host delete [options] [args]
		
  Delete a drago host.
  
  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *HostDeleteCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	c.UI.Output("Not implemented\n")

	return 0
}
