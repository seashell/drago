package command

import (
	"context"
	"strings"
)

type HostCreateCommand struct {
	BaseCommand
}

func (c *HostCreateCommand) Name() string {
	return "host create"
}

func (c *HostCreateCommand) Synopsis() string {
	return "Create a drago host"
}

func (c *HostCreateCommand) Help() string {
	help := `
Usage: drago host create [options] [args]
		
  Create a drago host.
  
  General Options:

  ` + GlobalOptions() + `

`
	return strings.TrimSpace(help)
}

func (c *HostCreateCommand) Run(ctx context.Context, args []string) int {

	flags := FlagSet(c.Name())

	flags.Usage = func() {
		c.UI.Output(c.Help())
	}

	var name string
	flags.StringVar(&name, "name", "", "")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate args
	args = flags.Args()
	if len(args) != 0 {
		c.UI.Error("This command takes no arguments")
		c.UI.Error(DefaultErrorMessage(c))
		return 1
	}

	c.UI.Output("Not implemented\n")

	return 0
}
