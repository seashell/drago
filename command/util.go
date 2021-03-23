package command

import (
	"fmt"

	api "github.com/seashell/drago/api"
	"github.com/seashell/drago/drago/structs"
)

// Returns the node ID of the local agent, in case it is a client.
// Otherwise, returns an error.
func localAgentNodeID(api *api.Client) (string, error) {

	self, err := api.Agent().Self(&structs.QueryOptions{})
	if err != nil {
		return "", fmt.Errorf("could not retrieve agent info: %s", err)
	}

	clientStats, ok := self.Stats["client"]
	if !ok {
		return "", fmt.Errorf("not running in client mode")
	}

	nodeID, ok := clientStats["node_id"]
	if !ok {
		return "", fmt.Errorf("could not determine node id")
	}

	return nodeID, nil
}
