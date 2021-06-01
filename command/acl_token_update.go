package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// ACLTokenUpdateCommand :
type ACLTokenUpdateCommand struct {
	UI cli.UI
	Command

	// parsed flags
	json       bool
	name       string
	policyType string
	policies   []string
}

func (c *ACLTokenUpdateCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.name, "name", "", "")
	flags.StringVar(&c.policyType, "type", "", "")
	flags.StringSliceVar(&c.policies, "policy", []string{}, "")
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ACLTokenUpdateCommand) Name() string {
	return "acl token update"
}

// Synopsis :
func (c *ACLTokenUpdateCommand) Synopsis() string {
	return "Update ACL token"
}

// Run :
func (c *ACLTokenUpdateCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes one argument: <token_id>")
		c.UI.Error(`For additional help, try 'drago acl token update --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	token, err := api.ACLTokens().Update(&structs.ACLToken{
		Name:     c.name,
		Type:     c.policyType,
		Policies: c.policies,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating ACL token: %s", err))
		return 1
	}

	c.UI.Output(c.formatToken(token))

	return 0
}

// Help :
func (c *ACLTokenUpdateCommand) Help() string {
	h := `
Usage: drago acl token update <token_id> [options]

  Update an ACL token.

General Options:
` + GlobalOptions() + `

ACL Token Update Options:

  --name=<name>
	Sets the name of the ACL token.

  --type=<type>
	Sets the type of the ACL token. Must be either "client" or "management". If not provided, defaults to "client".

  --policy=<policy>
	Specifies policies to associate with a client token. Can be specified multiple times.

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ACLTokenUpdateCommand) formatToken(token *structs.ACLToken) string {

	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")
	formatted := map[string]interface{}{
		"id":        token.ID,
		"name":      token.Name,
		"type":      token.Type,
		"secret":    token.Secret,
		"policies":  token.Policies,
		"createdAt": token.CreatedAt,
		"updatedAt": token.UpdatedAt,
	}
	if err := enc.Encode(formatted); err != nil {
		c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
	}

	s := b.String()

	if c.json {
		return s
	}

	return cleanJSONString(s)
}
