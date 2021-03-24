package command

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	table "github.com/rodaine/table"
	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
)

// NodeListCommand :
type NodeListCommand struct {
	UI cli.UI

	// Parsed flags
	self bool
	json bool

	Command
}

func (c *NodeListCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.self, "self", false, "")
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *NodeListCommand) Name() string {
	return "node status"
}

// Synopsis :
func (c *NodeListCommand) Synopsis() string {
	return "Display status information about nodes"
}

// Run :
func (c *NodeListCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 1 {
		c.UI.Error("This command takes either one or no arguments")
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	// Print status of multiple nodes
	if len(args) == 0 && !c.self {

		nodes, err := api.Nodes().List(&structs.QueryOptions{AuthToken: c.token})
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

	// Print status of a single node
	var nodeID string
	if !c.self {
		nodeID = args[0]
	} else {
		if nodeID, err = localAgentNodeID(api); err != nil {
			c.UI.Error(fmt.Sprintf("Error determining local node ID: %s", err))
			return 1
		}
	}

	node, err := api.Nodes().Get(nodeID, &structs.QueryOptions{AuthToken: c.token})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving node status: %s", err))
		return 1
	}

	c.UI.Output(c.formatNode(node))

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

  -self
    Query the status of the local node.

  -json=<bool>
    Enable JSON output.

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
