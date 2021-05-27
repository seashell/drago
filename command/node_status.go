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

// NodeStatusCommand :
type NodeStatusCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	self bool
	json bool
}

func (c *NodeStatusCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.self, "self", false, "")
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *NodeStatusCommand) Name() string {
	return "node status"
}

// Synopsis :
func (c *NodeStatusCommand) Synopsis() string {
	return "Display status of existing nodes"
}

// Run :
func (c *NodeStatusCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes either one or no arguments")
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	if len(args) == 0 && !c.self {

		filters := map[string][]string{}

		// Print status of multiple nodes
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

	node, err := api.Nodes().Get(nodeID)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving node status: %s", err))
		return 1
	}

	c.UI.Output(c.formatNode(node))

	return 0
}

// Help :
func (c *NodeStatusCommand) Help() string {
	h := `
Usage: drago node status <node_id> [options]

  Display node status information.

  If a node ID is passed, information for that specific node will be displayed.
  If no node ID's are passed, then a short-hand list of all nodes will be displayed.
  The -self flag is useful to quickly access the status of the local node.

  If ACLs are enabled, this option requires a token with the 'node:read' capability.

General Options:
` + GlobalOptions() + `

Node Status Options:

  --self
    Query the status of the local node.

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *NodeStatusCommand) formatNodeList(nodes []*structs.NodeListStub) string {

	var b bytes.Buffer
	fnodes := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, node := range nodes {
			fnodes = append(fnodes, map[string]string{
				"ID":     node.ID,
				"Status": node.Status,
			})
		}
		if err := enc.Encode(fnodes); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("NODE ID", "STATUS").WithWriter(&b)
		for _, node := range nodes {
			tbl.AddRow(node.ID, node.Status)
		}
		tbl.Print()
	}

	return b.String()
}

func (c *NodeStatusCommand) formatNode(node *structs.Node) string {

	var b bytes.Buffer

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")

		fnode := map[string]string{
			"ID":     node.ID,
			"Status": node.Status,
		}

		if err := enc.Encode(fnode); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}

	} else {
		tbl := table.New("NODE ID", "STATUS").WithWriter(&b)
		tbl.AddRow(node.ID, node.Status)
		tbl.Print()
	}

	return b.String()
}
