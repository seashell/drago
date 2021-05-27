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

// NetworkListCommand :
type NetworkListCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json bool
}

func (c *NetworkListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *NetworkListCommand) Name() string {
	return "network list"
}

// Synopsis :
func (c *NetworkListCommand) Synopsis() string {
	return "Display a list of networks"
}

// Run :
func (c *NetworkListCommand) Run(ctx context.Context, args []string) int {

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

	networks, err := api.Networks().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving networks: %s", err))
		return 1
	}

	if len(networks) == 0 {
		return 0
	}

	c.UI.Output(c.formatNetworkList(networks))

	return 0
}

// Help :
func (c *NetworkListCommand) Help() string {
	h := `
Usage: drago network list [options]

  Lists networks managed by Drago.

  If ACLs are enabled, this option requires a token with the 'network:read' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *NetworkListCommand) formatNetworkList(networks []*structs.NetworkListStub) string {

	var b bytes.Buffer
	fnetworks := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, network := range networks {
			fnetworks = append(fnetworks, map[string]string{
				"ID":           network.ID,
				"Name":         network.Name,
				"AddressRange": network.AddressRange,
			})
		}
		if err := enc.Encode(fnetworks); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("NETWORK ID", "NAME", "ADDRESS RANGE").WithWriter(&b)
		for _, network := range networks {
			tbl.AddRow(network.ID, network.Name, network.AddressRange)
		}
		tbl.Print()
	}

	return b.String()
}

func (c *NetworkListCommand) formatNetwork(network *structs.Network) string {

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
