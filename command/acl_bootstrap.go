package command

import (
	"context"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

// ACLBootstrapCommand :
type ACLBootstrapCommand struct {
	UI cli.UI
}

// Name :
func (c *ACLBootstrapCommand) Name() string {
	return "acl bootstrap"
}

// Synopsis :
func (c *ACLBootstrapCommand) Synopsis() string {
	return "Bootstrap the ACL system"
}

// Run :
func (c *ACLBootstrapCommand) Run(ctx context.Context, args []string) int {
	c.UI.Warn("Command not implemented")
	return 0
}

// Help :
func (c *ACLBootstrapCommand) Help() string {
	h := `
Usage: drago acl bootstrap [options]

  Bootstrap is used to bootstrap the ACL system and get an initial token.

General Options:
` + GlobalOptions() + `
`
	return strings.TrimSpace(h)
}
