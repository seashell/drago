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

// ACLPolicyListCommand :
type ACLPolicyListCommand struct {
	UI cli.UI

	// Parsed flags
	json bool

	Command
}

func (c *ACLPolicyListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ACLPolicyListCommand) Name() string {
	return "acl policy list"
}

// Synopsis :
func (c *ACLPolicyListCommand) Synopsis() string {
	return "List ACL policies"
}

// Run :
func (c *ACLPolicyListCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes no arguments")
		c.UI.Error(`For additional help, try 'drago acl policy list --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	policies, err := api.ACLPolicies().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving ACL policies: %s", err))
		return 1
	}

	if len(policies) == 0 {
		return 0
	}

	c.UI.Output(c.formatPolicyList(policies))

	return 0
}

// Help :
func (c *ACLPolicyListCommand) Help() string {
	h := `
Usage: drago acl policy list [options]

  List existing ACL policies.

General Options:
` + GlobalOptions() + `

ACL Policy List Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ACLPolicyListCommand) formatPolicyList(policies []*structs.ACLPolicyListStub) string {

	var b bytes.Buffer
	fpolicies := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, policy := range policies {
			fpolicies = append(fpolicies, map[string]string{
				"Name":        policy.Name,
				"Description": policy.Description,
			})
		}
		if err := enc.Encode(fpolicies); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("POLICY NAME", "DESCRIPTION").WithWriter(&b)
		for _, policy := range policies {
			tbl.AddRow(policy.Name, policy.Description)
		}
		tbl.Print()
	}

	return b.String()
}
