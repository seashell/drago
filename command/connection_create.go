package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
	cli "github.com/seashell/drago/pkg/cli"
	"github.com/spf13/pflag"
)

// ACLTokenCreateCommand :
type ConnectionCreateCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json                bool
	persistentKeepalive int
	allowAll            bool
}

func (c *ConnectionCreateCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")
	flags.IntVar(&c.persistentKeepalive, "keepalive", 0, "")
	flags.BoolVar(&c.allowAll, "allow-all", false, "")

	return flags
}

// Name :
func (c *ConnectionCreateCommand) Name() string {
	return "connection create"
}

// Synopsis :
func (c *ConnectionCreateCommand) Synopsis() string {
	return "Create a new connection"
}

// Run :
func (c *ConnectionCreateCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 3 {
		c.UI.Error("This command takes three arguments: <src_node_id> <dst_node_id> <network>")
		c.UI.Error(`For additional help, try 'drago acl token create --help'`)
		return 1
	}

	networkID := "" // Get from Networks API based on Network.Name

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	connection, err := api.Connections().Create(&structs.Connection{
		NetworkID:           networkID,
		PersistentKeepalive: &c.persistentKeepalive,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating connection: %s", err))
		return 1
	}

	c.UI.Output(c.formatConnection(connection))

	return 0
}

// Help :
func (c *ConnectionCreateCommand) Help() string {
	h := `
Usage: drago connection create <src_node_id> <dst_node_id> <network> [options]

  Create is used to create a new connection between two nodes that have interfaces for the same network.

  If ACLs are enabled, this option requires a token with the 'connection:write' capability.

General Options:
` + GlobalOptions() + `

ACL Token Create Options:

  --json
    Enable JSON output.

  --allow-all
    Enables routing of all traffic in this connection.

  --keepalive=<seconds>
    Time interval between persistent keepalive packets. Defaults to 0, which disables the feature.

`
	return strings.TrimSpace(h)
}

func (c *ConnectionCreateCommand) formatConnection(connection *structs.Connection) string {

	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")
	formatted := map[string]interface{}{
		"id":        connection.ID,
		"networkId": connection.NetworkID,
	}
	if err := enc.Encode(formatted); err != nil {
		c.UI.Error(fmt.Sprintf("Error formatting JSON output: %s", err))
	}

	s := b.String()

	if c.json {
		return s
	}

	return cleanJSONString(s)
}
