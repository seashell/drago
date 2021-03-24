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

// InterfaceListCommand :
type InterfaceListCommand struct {
	UI cli.UI

	// Parsed flags
	json bool

	Command
}

func (c *InterfaceListCommand) FlagSet() *flag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

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

	ifaces, err := api.Interfaces().List(&structs.QueryOptions{AuthToken: c.token})
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

  -json=<bool>
    Enable JSON output.

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

func (c *InterfaceListCommand) formatInterface(iface *structs.Interface) string {

	var b bytes.Buffer

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")

		finterface := map[string]string{
			"ID":      iface.ID,
			"Name":    valueOrPlaceholder(iface.Name, "N/A"),
			"Address": valueOrPlaceholder(iface.Address, "N/A"),
		}

		if err := enc.Encode(finterface); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}

	} else {
		tbl := table.New("INTERFACE ID", "NAME", "ADDRESS RANGE").WithWriter(&b)
		tbl.AddRow(iface.ID, valueOrPlaceholder(iface.Name, "N/A"), valueOrPlaceholder(iface.Address, "N/A"))
		tbl.Print()
	}

	return b.String()
}
