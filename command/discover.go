package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// DiscoverCommand :
type DiscoverCommand struct {
	UI cli.UI
}

// Name :
func (c *DiscoverCommand) Name() string {
	return "discover"
}

// Synopsis :
func (c *DiscoverCommand) Synopsis() string {
	return "Discover IPs based on metadata"
}

// Run :
func (c *DiscoverCommand) Run(ctx context.Context, args []string) int {

	return 0
}

// Help :
func (c *DiscoverCommand) Help() string {
	h := `
Usage: drago discover <label> [options]
	
  Retrieve the IP addresses of all interfaces belonging
  to nodes whose metadata contains an entry identified by
  the provided label.

General Options:
` + GlobalOptions() + `

Discover Options:

  --host=<address>
    The address of the host to be queried.
	
`
	return strings.TrimSpace(h)
}
