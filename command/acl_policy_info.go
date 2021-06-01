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

// ACLPolicyInfoCommand :
type ACLPolicyInfoCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json bool
}

func (c *ACLPolicyInfoCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ACLPolicyInfoCommand) Name() string {
	return "acl policy info"
}

// Synopsis :
func (c *ACLPolicyInfoCommand) Synopsis() string {
	return "Display details about an existing ACL policy"
}

// Run :
func (c *ACLPolicyInfoCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <name>")
		c.UI.Error(`For additional help, try 'drago acl policy info --help'`)
		return 1
	}

	name := args[0]

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	policy, err := api.ACLPolicies().Get(name)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving ACL policy: %s", err))
		return 1
	}

	c.UI.Output(c.formatPolicy(policy))

	return 0
}

// Help :
func (c *ACLPolicyInfoCommand) Help() string {
	h := `
Usage: drago acl policy info <name> [options]

  Display information on an existing ACL policy.

  Use the --json flag to see a detailed list of the rules associated with the policy.

General Options:
` + GlobalOptions() + `

ACL Policy Info Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ACLPolicyInfoCommand) formatPolicy(policy *structs.ACLPolicy) string {

	var b bytes.Buffer

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		fpolicy := map[string]interface{}{
			"name":        policy.Name,
			"description": policy.Description,
			"rules":       policy.Rules,
		}
		if err := enc.Encode(fpolicy); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("POLICY NAME", "DESCRIPTION", "RULES").WithWriter(&b)
		tbl.AddRow(policy.Name, policy.Description, len(policy.Rules))
		tbl.Print()
	}

	return b.String()
}
