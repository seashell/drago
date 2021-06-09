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

// ConnectionUpdateRulesCommand :
type ConnectionUpdateRulesCommand struct {
	UI cli.UI
	Command

	// Parsed flags
	allow     []string
	allowAll  bool
	allowNone bool
	on        string
}

func (c *ConnectionUpdateRulesCommand) FlagSet() *pflag.FlagSet {

	flags := c.Command.FlagSet(c.Name())

	flags.Usage = func() { c.UI.Output("\n" + c.Help() + "\n") }

	// General options
	flags.StringVar(&c.on, "on", "", "")
	flags.StringSliceVar(&c.allow, "allow", []string{}, "")
	flags.BoolVar(&c.allowAll, "allow-all", false, "")
	flags.BoolVar(&c.allowNone, "allow-none", false, "")

	return flags
}

// Name :
func (c *ConnectionUpdateRulesCommand) Name() string {
	return "connection update rules"
}

// Synopsis :
func (c *ConnectionUpdateRulesCommand) Synopsis() string {
	return "Update routing rules of a connection"
}

// Run :
func (c *ConnectionUpdateRulesCommand) Run(ctx context.Context, args []string) int {

	flags := c.FlagSet()

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 3 {
		c.UI.Error("This command takes one argument: <connection_id>")
		c.UI.Error(`For additional help, try 'drago connection update-rules --help'`)
		return 1
	}

	connectionID := ""

	// Get the HTTP client
	api, err := c.Command.APIClient()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting up API client: %s", err))
		return 1
	}

	conn, err := api.Connections().Get(connectionID)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting connection: %s", err))
		return 1
	}

	if c.allowAll {
		network, err := api.Networks().Get(conn.NetworkID)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting network: %s", err))
			return 1
		}

		if c.on == "" {
			for i := range conn.PeerSettings {
				conn.PeerSettings[i].RoutingRules.AllowedIPs = []string{network.AddressRange}
			}
		} else {
			for i := range conn.PeerSettings {
				if conn.PeerSettings[i].InterfaceID == c.on || conn.PeerSettings[i].NodeID == c.on {
					conn.PeerSettings[i].RoutingRules.AllowedIPs = []string{network.AddressRange}

					break
				}
			}
		}
	} else {
		if c.on == "" {
			for i := range conn.PeerSettings {
				conn.PeerSettings[i].RoutingRules.AllowedIPs = c.allow
			}
		} else {
			for i := range conn.PeerSettings {
				if conn.PeerSettings[i].InterfaceID == c.on || conn.PeerSettings[i].NodeID == c.on {
					conn.PeerSettings[i].RoutingRules.AllowedIPs = c.allow

					break
				}
			}
		}
	}

	err = api.Connections().Update(conn)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error updating connection: %s", err))
		return 1
	}

	return 0
}

// Help :
func (c *ConnectionUpdateRulesCommand) Help() string {
	h := `
Usage: drago connection update-rules <connection_id> [options]

  Update is used to update the routing rules enforced on each interface of a connection.

  If ACLs are enabled, this option requires a token with the 'connection:write' capability.

General Options:
` + GlobalOptions() + `

ACL Token Create Options:

  --allow
    Allow routing traffic to this address by the specified end of the connection.

  --allow-all
    Enable routing of all traffic by the specified end of the connection.

  --allow-none
    Disable routing of all traffic by the specified end of the connection.

  --on=<id>
    Node or interface ID specifying to which end of the connection the rules should be applied.

`
	return strings.TrimSpace(h)
}

func (c *ConnectionUpdateRulesCommand) formatConnection(connection *structs.Connection) string {

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

	return cleanJSONString(s)
}
