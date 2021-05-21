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

	// Parsed flags
	json    bool
	node    []string
	network []string

	Command
}

func (c *InterfaceListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")
	flags.StringSliceVar(&c.node, "node", []string{}, "")
	flags.StringSliceVar(&c.network, "network", []string{}, "")

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
		return 1
	}

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	filters := map[string][]string{}

	if len(c.network) > 0 {
		filters["network"] = c.network
	}

	if len(c.node) > 0 {
		filters["node"] = c.node
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

  Lists interfaces managed by Drago.

  If ACLs are enabled, this option requires a token with the 'interface:read' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  --json
    Enable JSON output.

  --node=<id>
    Filter results by node ID.

  --network=<id>
    Filter results by network ID.

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
				"ID":      iface.ID,
				"Name":    valueOrPlaceholder(iface.Name, "N/A"),
				"Address": valueOrPlaceholder(iface.Address, "N/A"),
			})
		}
		if err := enc.Encode(fifaces); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("INTERFACE ID", "NAME", "ADDRESS").WithWriter(&b)
		for _, iface := range interfaces {
			tbl.AddRow(iface.ID, valueOrPlaceholder(iface.Name, "N/A"), valueOrPlaceholder(iface.Address, "N/A"))
		}
		tbl.Print()
	}

	return b.String()
}
