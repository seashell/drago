package command

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
)

// NodeStatusCommand :
type NodeStatusCommand struct {
	UI cli.UI

	// Parsed flags
	self bool

	Command
}

func (c *NodeStatusCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options (available in both client and server modes)
	flags.BoolVar(&c.self, "self", false, "")

	return flags
}

// Name :
func (c *NodeStatusCommand) Name() string {
	return "node status"
}

// Synopsis :
func (c *NodeStatusCommand) Synopsis() string {
	return "Display status information about nodes"
}

// Run :
func (c *NodeStatusCommand) Run(ctx context.Context, args []string) int {

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

		for _, node := range nodes {
			c.UI.Output(nodeStubStr(node))
		}

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

	c.UI.Output(nodeStr(node))

	return 0
}

// Help :
func (c *NodeStatusCommand) Help() string {
	h := `
Usage: drago node status [options] <node>

  Display node status information.

  If a node ID is passed, information for that specific node will be displayed.
  If no node ID's are passed, then a short-hand list of all nodes will be displayed.
  The -self flag is useful to quickly access the status of the local node.

  If ACLs are enabled, this option requires a token with the 'node:read' capability.

General Options:
` + GlobalOptions() + `

Node Status Options:

  -self
    Query the status of the local node.

 `
	return strings.TrimSpace(h)
}

func nodeStubStr(node *structs.NodeListStub) string {
	return fmt.Sprintf("%s    %s", node.ID, node.Status)
}

func nodeStr(node *structs.Node) string {
	return fmt.Sprintf("%s   %s", node.ID, node.Status)
}
