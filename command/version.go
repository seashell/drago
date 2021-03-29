package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
	version "github.com/seashell/drago/version"
)

// VersionCommand :
type VersionCommand struct {
	UI cli.UI
}

// Name :
func (c *VersionCommand) Name() string {
	return "version"
}

// Synopsis :
func (c *VersionCommand) Synopsis() string {
	return "Print the Drago version"
}

// Run :
func (c *VersionCommand) Run(ctx context.Context, args []string) int {
	c.UI.Output(version.GetVersion().VersionNumber())
	return 0
}

// Help :
func (c *VersionCommand) Help() string {
	h := `
Usage: drago version [options]

  Version prints out the Drago version.

General Options:
` + GlobalOptions() + `
`
	return strings.TrimSpace(h)
}
