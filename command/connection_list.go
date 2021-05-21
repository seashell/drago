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

// ConnectionListCommand :
type ConnectionListCommand struct {
	UI cli.UI

	// Parsed flags
	json bool

	Command
}

func (c *ConnectionListCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")

	return flags
}

// Name :
func (c *ConnectionListCommand) Name() string {
	return "connection list"
}

// Synopsis :
func (c *ConnectionListCommand) Synopsis() string {
	return "Display a list of connections"
}

// Run :
func (c *ConnectionListCommand) Run(ctx context.Context, args []string) int {

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

	connections, err := api.Connections().List()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving connections: %s", err))
		return 1
	}

	if len(connections) == 0 {
		return 0
	}

	c.UI.Output(c.formatConnectionList(connections))

	return 0
}

// Help :
func (c *ConnectionListCommand) Help() string {
	h := `
Usage: drago connection list [options]

  Lists connections managed by Drago.

  If ACLs are enabled, this option requires a token with the 'connection:read' capability.

General Options:
` + GlobalOptions() + `

Connection List Options:

  --json
    Enable JSON output.

`
	return strings.TrimSpace(h)
}

func (c *ConnectionListCommand) formatConnectionList(connections []*structs.ConnectionListStub) string {

	var b bytes.Buffer
	fconnections := []interface{}{}

	if c.json {
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "    ")
		for _, conn := range connections {
			fconnections = append(fconnections, map[string]string{
				"ID": conn.ID,
			})
		}
		if err := enc.Encode(fconnections); err != nil {
			c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
		}
	} else {
		tbl := table.New("CONNECTION ID").WithWriter(&b)
		for _, conn := range connections {
			tbl.AddRow(conn.ID)
		}
		tbl.Print()
	}

	return b.String()
}
