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

// ACLTokenCreateCommand :
type ACLTokenCreateCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json       bool
	name       string
	policyType string
	policies   []string
}

func (c *ACLTokenCreateCommand) FlagSet() *pflag.FlagSet {

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
func (c *ACLTokenCreateCommand) Name() string {
	return "acl token create"
}

// Synopsis :
func (c *ACLTokenCreateCommand) Synopsis() string {
	return "Create a new ACL token"
}

// Run :
func (c *ACLTokenCreateCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes no arguments")
		c.UI.Error(`For additional help, try 'drago acl token create --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	token, err := api.ACLTokens().Create(&structs.ACLToken{
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
func (c *ACLTokenCreateCommand) Help() string {
	h := `
Usage: drago acl token create [options]

  Create is used to issue a new ACL token. Requires a management token.

General Options:
` + GlobalOptions() + `

ACL Token Create Options:

  --name=<name>
    Sets the human readable name for the ACL token.

  --type=<type>
    Sets the type of token. Must be one of "client" (default), or "management".

  --policy=<policy>
    Specifies a policy to associate with client tokens.

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ACLTokenCreateCommand) formatToken(token *structs.ACLToken) string {

	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")
	formatted := map[string]interface{}{
		"ID":        token.ID,
		"Name":      token.Name,
		"Type":      token.Type,
		"Secret":    token.Secret,
		"Policies":  token.Policies,
		"CreatedAt": token.CreatedAt,
		"UpdatedAt": token.UpdatedAt,
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
