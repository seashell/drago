package acl

import (
	"context"
	"fmt"
	"path"
	"strings"

	radix "github.com/seashell/drago/pkg/radix"
)

const (
	capabilityDeny = "deny"
)

// ACL is used to convert a set of policies into a structure that
// can be efficiently evaluated to determine if an action is allowed.
type ACL struct {
	privileged   bool
	capabilities map[string]*radix.Tree
}

// CheckAuthorized verifies whether the ACL is authorized to perform a specific action.
// If the ACL is not authorized, an error is returned, which provides more details.
// If an operation is not explicitly enabled in the ACL, it is forbidden by default.
func (a *ACL) CheckAuthorized(ctx context.Context, res string, path string, op string) error {
	// A privileged ACL is able to do anything.
	if a.privileged {
		return nil
	}

	capabilities, err := a.queryCapabilities(ctx, res, path)
	if err != nil {
		return err
	}

	if capabilities.HasCapability(op) {
		return nil
	}

	return ErrNotAuthorized
}

// queryCapabilities searches the ACL for all rules matching the queried resource
// instance, and merges the capabilities they enable into one single capabilityMap,
// which can then be used to easily verify whether a capability is set or not.
func (a *ACL) queryCapabilities(ctx context.Context, resource string, instance string) (capMap, error) {

	capTree, ok := a.capabilities[resource]
	if !ok {
		return nil, ErrInvalidResource
	}

	if v, found := capTree.Get(instance); found {
		return v.(capMap), nil
	}

	// Find all patterns in the capability tree
	// matching the queried instance.
	merged := make(capMap)
	capTree.Walk(func(pattern string, raw interface{}) bool {
		if matched, _ := path.Match(pattern, instance); matched {
			for cap := range raw.(capMap) {
				merged.AddCapability(cap)
			}
			return false
		}
		return true
	})

	return merged, nil
}

// String builds a string representation of an ACL which
// can be useful for debugging purposes.
func (a *ACL) String() string {
	s := ""
	s += fmt.Sprintf("- management = %v\n", a.privileged)
	for res, capTree := range a.capabilities {
		s += fmt.Sprintf("-- %s capabilities\n", res)
		capTree.Walk(func(k string, v interface{}) bool {
			s += fmt.Sprintf("    %s%s --> ", k, strings.Repeat(" ", 18-len(k)))
			for c := range v.(capMap) {
				s += fmt.Sprintf("%s  ", c)
			}
			s += fmt.Sprintf("\n")
			return false
		})
	}
	return s
}

// capMap is a map meant for making it easy/efficient to
// evaluate whether or not a given capability is enabled. If an
// entry exists in the map for a specific capability, it is enabled.
// Otherwise, it is not.
type capMap map[string]struct{}

func (m capMap) AddCapability(cap string) {
	m[cap] = struct{}{}
}

func (m capMap) RemoveCapability(cap string) {
	delete(m, cap)
}

func (m capMap) HasCapability(cap string) bool {
	_, found := m[cap]
	return found
}

func (m capMap) Clear() {
	for cap := range m {
		delete(m, cap)
	}
}
