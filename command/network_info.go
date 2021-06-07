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

// NetworkInfoCommand :
type NetworkInfoCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json bool
}

func (c *NetworkInfoCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
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
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <network>")
		c.UI.Error(`For additional help, try 'drago network info --help'`)
		return 1
	}

	name := args[0]
	id := ""

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	networks, err := api.Networks().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting networks: %s", err))
		return 1
	}

	for _, n := range networks {
		if n.Name == name {
			id = n.ID

			break
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
  Usage: drago network info <network> [options]

  Display detailed information about an existing Drago network.

  If ACLs are enabled, this option requires a token with the 'network:read' capability.

General Options:
` + GlobalOptions() + `

Network Info Options:

  --json
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
			"id":           network.ID,
			"name":         network.Name,
			"addressRange": network.AddressRange,
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
