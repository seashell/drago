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

// InterfaceListCommand :
type InterfaceListCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json    bool
	self    bool
	node    string
	network string
}

func (c *InterfaceListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")
	flags.BoolVar(&c.self, "self", false, "")
	flags.StringVar(&c.node, "node", "", "")
	flags.StringVar(&c.network, "network", "", "")

	return flags
}

// Name :
func (c *InterfaceListCommand) Name() string {
	return "interface list"
}

// Synopsis :
func (c *InterfaceListCommand) Synopsis() string {
	return "Display a list of interfaces"
}

// Run :
func (c *InterfaceListCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 0 {
		c.UI.Error("This command takes no arguments")
		c.UI.Error(`For additional help, try 'drago interface list --help'`)
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	filters := map[string][]string{}
	networkID := ""

	if len(c.network) > 0 {
		networks, err := api.Networks().List()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error retrieving networks: %s", err))
			return 1
		}

		for _, network := range networks {
			if c.network == network.Name {
				networkID = network.ID

				break
			}
		}
	}

	if c.self && len(c.node) > 0 {
		c.UI.Error("Can not have both the --self and the --node flags.")
		return 1
	}

	if len(c.network) > 0 {
		filters["network"] = []string{networkID}
	}

	if len(c.node) > 0 {
		filters["node"] = []string{c.node}
	}

	if c.self {
		nodeID := ""

		if nodeID, err = localAgentNodeID(api); err != nil {
			c.UI.Error(fmt.Sprintf("Error determining local node ID: %s", err))
			return 1
		}

		filters["node"] = []string{nodeID}
	}

	ifaces, err := api.Interfaces().List(filters)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving interfaces: %s", err))
		return 1
	}

	if len(ifaces) == 0 {
		return 0
	}

	c.UI.Output(c.formatInterfaceList(ifaces))

	return 0
}

// Help :
func (c *InterfaceListCommand) Help() string {
	h := `
Usage: drago interface list [options]

  List interfaces managed by Drago.

  If ACLs are enabled, this option requires a token with the 'interface:read' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  --json
    Enable JSON output.

  --self
    Filter results by the local node ID. Can not be used with the --node filter flag.

  --node=<node_id>
    Filter results by node ID. Can not be used with the --self filter flag.

  --network=<network>
    Filter results by network.

`
	return strings.TrimSpace(h)
}

func (c *InterfaceListCommand) formatInterfaceList(interfaces []*structs.InterfaceListStub) string {

	var b bytes.Buffer
	fifaces := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, iface := range interfaces {
			fifaces = append(fifaces, map[string]string{
				"id":      iface.ID,
				"address": valueOrPlaceholder(iface.Address, "N/A"),
				"network": iface.NetworkID,
				"node":    iface.NodeID,
			})
		}
		if err := enc.Encode(fifaces); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("INTERFACE ID", "ADDRESS", "NETWORK ID", "NODE ID").WithWriter(&b)
		for _, iface := range interfaces {
			tbl.AddRow(iface.ID, valueOrPlaceholder(iface.Address, "N/A"), iface.NetworkID, iface.NodeID)
		}
		tbl.Print()
	}

	return b.String()
}
