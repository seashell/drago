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

// NodeListCommand :
type NodeListCommand struct {
	UI cli.UI

	// Parsed flags
	json   bool
	status string
	meta   []string

	Command
}

func (c *NodeListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")
	flags.StringVar(&c.status, "status", "*", "")
	flags.StringSliceVar(&c.meta, "meta", []string{}, "")

	return flags
}

// Name :
func (c *NodeListCommand) Name() string {
	return "node status"
}

// Synopsis :
func (c *NodeListCommand) Synopsis() string {
	return "List existing nodes"
}

// Run :
func (c *NodeListCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 0 {
		c.UI.Error("This command takes no arguments")
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	filters := map[string][]string{}
	filters["meta"] = c.meta
	filters["status"] = []string{c.status}

	nodes, err := api.Nodes().List(filters)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving node status: %s", err))
		return 1
	}

	if len(nodes) == 0 {
		return 0
	}

	c.UI.Output(c.formatNodeList(nodes))

	return 0

}

// Help :
func (c *NodeListCommand) Help() string {
	h := `
Usage: drago node list [options]

  List nodes registered on Drago.

  If ACLs are enabled, this option requires a token with the 'node:read' capability.

General Options:
` + GlobalOptions() + `

Node List Options:

  --json=<bool>
    Enable JSON output.
  
  --meta=<key:value>
    Filter nodes by metadata.

  --status=<initializing|ready|down>
    Filter nodes by status.

`
	return strings.TrimSpace(h)
}

func (c *NodeListCommand) formatNodeList(nodes []*structs.NodeListStub) string {

	var b bytes.Buffer
	fnodes := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, node := range nodes {
			fnodes = append(fnodes, map[string]string{
				"ID":     node.ID,
				"Name":   node.Name,
				"Status": node.Status,
			})
		}
		if err := enc.Encode(fnodes); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("NODE ID", "NAME", "STATUS").WithWriter(&b)
		for _, node := range nodes {
			tbl.AddRow(node.ID, node.Name, node.Status)
		}
		tbl.Print()
	}

	return b.String()
}

func (c *NodeListCommand) formatNode(node *structs.Node) string {

	var b bytes.Buffer

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")

		fnode := map[string]string{
			"ID":     node.ID,
			"Name":   node.Name,
			"Status": node.Status,
		}

		if err := enc.Encode(fnode); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}

	} else {
		tbl := table.New("NODE ID", "NAME", "STATUS").WithWriter(&b)
		tbl.AddRow(node.ID, node.Name, node.Status)
		tbl.Print()
	}

	return b.String()
}
