package command

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
)

// ACLTokenCreateCommand :
type ACLTokenCreateCommand struct {
	UI cli.UI

	// Parsed flags
	tname     string
	ttype     string
	tpolicies manyStrings

	Command
}

func (c *ACLTokenCreateCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.tname, "name", "", "")
	flags.StringVar(&c.ttype, "type", "", "")
	flags.Var(&c.tpolicies, "type", "")

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
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	t := &structs.ACLToken{
		Name:     c.tname,
		Type:     c.ttype,
		Policies: c.tpolicies,
	}

	token, err := api.ACLTokens().Create(t, &structs.QueryOptions{AuthToken: c.token})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving ACL token: %s", err))
		return 1
	}

	c.UI.Output(c.formatToken(token))

	return 0
}

// Help :
func (c *ACLTokenCreateCommand) Help() string {
	h := `
Usage: drago acl token info <name> [options]

  Display information on an existing ACL token.

  Use the -json flag to see a detailed list of the rules associated with the token.

General Options:
` + GlobalOptions() + `

ACL Token Create Options:

  -name=""
    Sets the human readable name for the ACL token.

  -type="client"
    Sets the type of token. Must be one of "client" (default), or "management".

  -policy="string"
    Specifies a policy to associate with client tokens.

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

	return b.String()
}
