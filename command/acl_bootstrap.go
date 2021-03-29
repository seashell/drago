package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
)

// ACLBootstrapCommand :
type ACLBootstrapCommand struct {
	UI cli.UI

	Command
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

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	token, err := api.ACL().Bootstrap()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error bootstrapping ACL system: %s", err))
		return 1
	}

	c.UI.Output(c.formatACLToken(token))

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

func (c *ACLBootstrapCommand) formatACLToken(t *structs.ACLToken) string {

	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")

	if err := enc.Encode(t); err != nil {
		c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
	}

	return b.String()
}
