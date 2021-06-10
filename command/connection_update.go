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

// ConnectionUpdateCommand :
type ConnectionUpdateCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	json                bool
	persistentKeepalive int
	allowAll            bool
}

func (c *ConnectionUpdateCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())
	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.BoolVar(&c.json, "json", false, "")
	flags.IntVar(&c.persistentKeepalive, "keepalive", 0, "")
	flags.BoolVar(&c.allowAll, "allow-all", false, "")

	return flags
}

// Name :
func (c *ConnectionUpdateCommand) Name() string {
	return "connection update"
}

// Synopsis :
func (c *ConnectionUpdateCommand) Synopsis() string {
	return "Update an existing connection"
}

// Run :
func (c *ConnectionUpdateCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 3 {
		c.UI.Error("This command takes one argument: <connection_id>")
		c.UI.Error(`For additional help, try 'drago connection update --help'`)
		return 1
	}

	connectionID := ""

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	conn := &structs.Connection{
		ID:                  connectionID,
		PersistentKeepalive: &c.persistentKeepalive,
	}

	rcvConn, err := api.Connections().Update(conn)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error updating connection: %s", err))
		return 1
	}

	c.UI.Output(c.formatConnection(rcvConn))

	return 0
}

// Help :
func (c *ConnectionUpdateCommand) Help() string {
	h := `
Usage: drago connection update <connection_id> [options]

  Update is used to update an existing connection between two nodes that have interfaces in the same network.

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

func (c *ConnectionUpdateCommand) formatConnection(connection *structs.Connection) string {

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
