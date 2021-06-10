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

// InterfaceUpdateCommand :
type InterfaceUpdateCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	address string
	json    bool
}

func (c *InterfaceUpdateCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.address, "address", "", "")
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *InterfaceUpdateCommand) Name() string {
	return "interface update"
}

// Synopsis :
func (c *InterfaceUpdateCommand) Synopsis() string {
	return "Update an existing interface"
}

// Run :
func (c *InterfaceUpdateCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		c.UI.Error("This command takes one argument: <interface_id>")
		c.UI.Error(`For additional help, try 'drago interface update --help'`)
		return 1
	}

	id := args[0]

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	iface, err := api.Interfaces().Update(&structs.Interface{
		ID:      id,
		Address: &c.address,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error updating interface: %s", err))
		return 1
	}

	c.UI.Output(c.formatInterface(iface))

	return 0
}

// Help :
func (c *InterfaceUpdateCommand) Help() string {
	h := `
Usage: drago interface update <interface_id> [options]

  Update an existing interface.

  If ACLs are enabled, this option requires a token with the 'interface:write' capability.

General Options:
` + GlobalOptions() + `

Network List Options:

  --json
	Enable JSON output.

  --address=<addr>
    Interface address.

`
	return strings.TrimSpace(h)
}

func (c *InterfaceUpdateCommand) formatInterface(iface *structs.Interface) string {

	var b bytes.Buffer

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")

		fiface := map[string]string{
			"id":      iface.ID,
			"address": valueOrPlaceholder(iface.Address, "N/A"),
			"network": iface.NetworkID,
			"node":    iface.NodeID,
		}

		if err := enc.Encode(fiface); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}

	} else {
		tbl := table.New("INTERFACE ID", "ADDRESS", "NETWORK ID", "NODE ID").WithWriter(&b)
		tbl.AddRow(iface.ID, valueOrPlaceholder(iface.Address, "N/A"), iface.NetworkID, iface.NodeID)
		tbl.Print()
	}

	return b.String()
}
