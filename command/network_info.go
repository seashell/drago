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

// NetworkInfoCommand :
type NetworkInfoCommand struct {
	UI cli.UI

	// Parsed flags
	name string
	json bool

	Command
}

func (c *NetworkInfoCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.name, "name", "", "")
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *NetworkInfoCommand) Name() string {
	return "network info"
}

// Synopsis :
func (c *NetworkInfoCommand) Synopsis() string {
	return "Display info about a network"
}

// Run :
func (c *NetworkInfoCommand) Run(ctx context.Context, args []string) int {

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

	id := args[0]

	if c.name != "" {
		networks, err := api.Networks().List()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting networks: %s", err))
			return 1
		}

		for _, n := range networks {
			if n.Name == c.name {
				if id != "" && n.ID != id {
					c.UI.Error("Error: name and ID belong to different networks")
					return 1
				}
				id = n.ID
			}
		}
	}

	if id == "" {
		c.UI.Error("Error: network not found")
		return 1
	}

	network, err := api.Networks().Get(id)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving network: %s", err))
		return 1
	}

	c.UI.Output(c.formatNetwork(network))

	return 0
}

// Help :
func (c *NetworkInfoCommand) Help() string {
	h := `
Usage: drago network create [options]

  Create a new Drago network.

  If ACLs are enabled, this option requires a token with the 'network:write' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  -name=""
    Sets the human readable name for the network.

  -name=""
    Sets the address range of the network, in CIDR notation.

  -json=<bool>
    Enable JSON output.

 `
	return strings.TrimSpace(h)
}

func (c *NetworkInfoCommand) formatNetwork(network *structs.Network) string {

	var b bytes.Buffer

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")

		fnetwork := map[string]string{
			"ID":           network.ID,
			"Name":         network.Name,
			"AddressRange": network.AddressRange,
		}

		if err := enc.Encode(fnetwork); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}

	} else {
		tbl := table.New("NETWORK ID", "NAME", "ADDRESS RANGE").WithWriter(&b)
		tbl.AddRow(network.ID, network.Name, network.AddressRange)
		tbl.Print()
	}

	return b.String()
}
