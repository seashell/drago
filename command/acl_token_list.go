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

// ACLTokenListCommand :
type ACLTokenListCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json bool
}

func (c *ACLTokenListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ACLTokenListCommand) Name() string {
	return "acl token list"
}

// Synopsis :
func (c *ACLTokenListCommand) Synopsis() string {
	return "List ACL tokens"
}

// Run :
func (c *ACLTokenListCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes no arguments")
		c.UI.Error(`For additional help, try 'drago acl token list --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	tokens, err := api.ACLTokens().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving ACL tokens: %s", err))
		return 1
	}

	if len(tokens) == 0 {
		return 0
	}

	c.UI.Output(c.formatTokenList(tokens))

	return 0
}

// Help :
func (c *ACLTokenListCommand) Help() string {
	h := `
Usage: drago acl token list [options]

  List existing ACL tokens.

General Options:
` + GlobalOptions() + `

ACL Policy List Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ACLTokenListCommand) formatTokenList(tokens []*structs.ACLTokenListStub) string {

	var b bytes.Buffer
	ftokens := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, token := range tokens {
			ftokens = append(ftokens, map[string]string{
				"ID":   token.ID,
				"Name": token.Name,
				"Type": token.Type,
			})
		}
		if err := enc.Encode(ftokens); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("TOKEN ID", "NAME", "TYPE").WithWriter(&b)
		for _, token := range tokens {
			tbl.AddRow(token.ID, token.Name, token.Type)
		}
		tbl.Print()
	}

	return b.String()
}
