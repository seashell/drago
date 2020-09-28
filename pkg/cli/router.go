package cli

import (
	"errors"
	"strings"
)

// Router ...
type Router struct {
	commands map[string]Command
}

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{
		commands: map[string]Command{},
	}
}

// AddCommand ...
func (r *Router) AddCommand(n string, cmd Command) error {
	r.commands[n] = cmd
	return nil
}

// GetCommand ...
func (r *Router) GetCommand(n string) (Command, error) {
	cmd, ok := r.commands[n]
	if !ok {
		return nil, errors.New("command not found ")
	}
	return cmd, nil
}

// AddMissingParents ensures that every registered command
// has a parent.
func (r *Router) AddMissingParents(genCmd func() Command) {

	missing := []string{}

	for key := range r.commands {

		idx := strings.LastIndex(key, " ")

		// Ignore top-level commands
		if idx == -1 {
			continue
		}

		key = key[:idx]

		// Ignore commands for which a parent already exists
		if _, ok := r.commands[key]; ok {
			continue
		}

		// If the command has no parent, update the map of missing commands
		missing = append(missing, key)
	}

	// Fill missing commands with mocks
	for _, key := range missing {
		r.AddCommand(key, genCmd())
	}

}

// GetSubcommands returns a map containing all commands
// which are below a given prefix
func (r *Router) GetSubcommands(prefix string) map[string]Command {
	// If prefix is not empty, make sure it ends in ' '
	if prefix != "" && prefix[len(prefix)-1] != ' ' {
		prefix += " "
	}

	var keys []string
	for k := range r.commands {
		if strings.HasPrefix(k, prefix) {
			if !strings.Contains(k[len(prefix):], " ") {
				// Ignore any sub-sub keys, i.e. "foo bar baz" when we want "foo bar"
				keys = append(keys, k)
			}
		}
	}

	// For each of the keys return that in the map
	result := make(map[string]Command, len(keys))
	for _, k := range keys {
		cmd, err := r.GetCommand(k)
		if err != nil {
			panic("not found: " + k)
		}
		result[k] = cmd
	}

	return result
}

// GetParent returns the parent of this subcommand, if there is one.
// Otherwise, "" is returned.
func (r *Router) GetParent(n string) string {
	if n == "" {
		return n
	}

	// Clear any trailing spaces and find the last space
	n = strings.TrimRight(n, " ")
	idx := strings.LastIndex(n, " ")

	if idx == -1 {
		return ""
	}

	return n[:idx]
}

// GetLongestPrefix return the longest prefix match among all commands
// registered in the router, considering the query string passed as argument
func (r *Router) GetLongestPrefix(s string) (string, interface{}, bool) {

	keys := make([]string, 0, len(r.commands))
	for k := range r.commands {
		keys = append(keys, k)
	}

	scores := make([]int, len(keys))
	for i, key := range keys {
		for j := range s {
			if j >= len(key) {
				break
			}
			if s[:j] != key[:j] {
				break
			}
			scores[i]++
		}
	}

	max := 0
	key := ""
	for i, score := range scores {
		if score > max {
			max = score
			key = keys[i]
		}
	}

	if key != "" {
		return key, r.commands[key], true
	}

	return key, nil, false
}
