package command

import (
	"fmt"
	"strings"

	api "github.com/seashell/drago/api"
)

// Returns the node ID of the local agent, in case it is a client.
// Otherwise, returns an error.
func localAgentNodeID(api *api.Client) (string, error) {

	self, err := api.Agent().Self()
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

func valueOrPlaceholder(s *string, p string) string {
	if s != nil {
		return *s
	}
	return p
}

// TODO: improve how we clean JSON strings
func cleanJSONString(s string) string {

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, `""`, "N/A")
	s = strings.ReplaceAll(s, `{}`, "\n      N/A")

	s = strings.ReplaceAll(s, "    ", "  ")
	s = strings.ReplaceAll(s, `,`, "")

	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, `"`, "")

	s = strings.TrimLeftFunc(s, func(r rune) bool {
		return r == '\n'
	})

	s = strings.TrimRightFunc(s, func(r rune) bool {
		return r == '\n' || r == ' '
	})

	return s

}
