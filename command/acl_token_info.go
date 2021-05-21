package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	table "github.com/rodaine/table"
	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// ACLTokenInfoCommand :
type ACLTokenInfoCommand struct {
	UI cli.UI

	// Parsed flags
	json bool

	Command
}

func (c *ACLTokenInfoCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ACLTokenInfoCommand) Name() string {
	return "acl token info"
}

// Synopsis :
func (c *ACLTokenInfoCommand) Synopsis() string {
	return "Display details about an existing ACL token"
}

// Run :
func (c *ACLTokenInfoCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.UI.Error("This command takes one argument: <id>")
		c.UI.Error(`For additional help, try 'drago acl token info --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	token, err := api.ACLTokens().Get(args[0])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving ACL token: %s", err))
		return 1
	}

	c.UI.Output(c.formatToken(token))

	return 0
}

// Help :
func (c *ACLTokenInfoCommand) Help() string {
	h := `
Usage: drago acl token info <token> [options]

  Display information on an existing ACL token.

General Options:
` + GlobalOptions() + `

ACL Token Info Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ACLTokenInfoCommand) formatToken(token *structs.ACLToken) string {

	var b bytes.Buffer

	if c.json {
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
	} else {
		tbl := table.New("TOKEN ID", "NAME", "TYPE", "SECRET", "POLICIES").WithWriter(&b)
		tbl.AddRow(token.ID, token.Name, token.Type, token.Secret, len(token.Policies))
		tbl.Print()
	}

	return b.String()
}
