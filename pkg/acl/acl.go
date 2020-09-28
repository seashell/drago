package acl

import (
	"fmt"
	"path"
	"sort"
	"strings"

	radix "github.com/seashell/drago/pkg/radix"
)

const (
	capabilityDeny = "deny"
)

// Global ACL configuration.
var config *Config

// ACL is used to convert a set of policies into a structure that
// can be efficiently evaluated to determine if an action is allowed.
type ACL struct {
	management   bool
	capabilities map[string]*radix.Tree
}

// Token ...
type Token interface {
	Policies() []string
}

// Policy represents a named set of rules for allowing
// operations on specific resources.
type Policy interface {
	Name() string
	Rules() []Rule
}

// Rule is used to allow operations on specifi resources. In
// addition to the name of the target resource type and the allowed
// operations, a rule also specifies an optional glob pattern for
// targeting specific instances of the resource.
type Rule interface {
	// Resource targeted by this rule.
	Resource() string
	// Pattern used to target specific instances
	// of the target resource, if applicable.
	Pattern() string
	// Capabilities contains the actions allowed on
	// instances of a resource matching this rule.
	Capabilities() []string
}

// Initialize initializes the ACL system based on user-provided
// configurations.
func Initialize(c *Config) {
	config = c
	config.LogActivity("acl successfully initialized")
}

// NewACL creates and properly initializes an ACL struct.
func NewACL() *ACL {
	return &ACL{
		capabilities: map[string]*radix.Tree{},
	}
}

// NewManagementACL creates an ACL with management capabilities.
func NewManagementACL() (*ACL, error) {
	return &ACL{management: true}, nil
}

// NewACLFromSecret creates an ACL from a secret.
func NewACLFromSecret(secret string) (*ACL, error) {

	// Resolve token by looking up its secret
	token, err := config.ResolveToken(secret)
	if token == nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Retrieve policy definitions from the repository.
	policies := []Policy{}
	for _, p := range token.Policies() {
		policy, err := config.ResolvePolicy(p)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}

	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Name() < policies[j].Name()
	})

	// Init ACL
	acl := NewACL()
	for resource := range config.Resources {
		acl.capabilities[resource] = radix.NewTree()
	}

	// Based on the token policies, we populate the
	// capability trees for each resource type.
	for _, policy := range policies {

		for _, res := range config.Resources {

			for _, rule := range policy.Rules() {

				var capabilities capabilityMap

				if rule.Resource() == res.Name {

					// Initialize capability map for this pattern if it has not yet been
					// initialized by a previously processed rule/policy.
					if leaf, found := acl.capabilities[res.Name].Get(rule.Pattern()); !found {
						capabilities = make(capabilityMap)
						acl.capabilities[res.Name].Set(rule.Pattern(), capabilities)
					} else {
						capabilities = leaf.(capabilityMap)
					}

					// Ignore all capabilities in the rule if a previously processed
					// rule/policy has explicitly denied access to resources matching this pattern.
					if capabilities.HasCapability(capabilityDeny) {
						break
					}

					// Expand aliases before processing capabilities.
					expanded := []string{}
					for _, cap := range rule.Capabilities() {
						if res.hasAlias(cap) {
							expanded = append(expanded, res.Aliases[cap].expand()...)
						} else {
							expanded = append(expanded, cap)
						}
					}

					// Set all capabilities in the rule.
					for _, cap := range expanded {
						if cap == capabilityDeny {
							capabilities.Clear()
							capabilities.AddCapability(capabilityDeny)
							break
						} else {
							capabilities.AddCapability(cap)
						}
					}
				}
			}
		}
	}
	return acl, nil
}

// IsAuthorized verifies whether the ACL is authorized to perform a specific action.
// If an operation is not explicitly enabled in the ACL, it is forbidden by default.
func (a *ACL) IsAuthorized(res string, inst string, op string) bool {

	// A management ACL is authorized to do anything.
	if a.management {
		return true
	}

	capabilities, err := a.queryCapabilities(res, inst)
	if err != nil {
		return false
	}

	return capabilities.HasCapability(op)
}

// queryCapabilities searches the ACL for all rules matching the queried resource
// instance, and merges the capabilities they enable into one single capabilityMap,
// which can then be used to easily verify whether a capability is set or not.
func (a *ACL) queryCapabilities(resource string, instance string) (capabilityMap, error) {

	capTree, ok := a.capabilities[resource]
	if !ok {
		return nil, fmt.Errorf("invalid resource type")
	}

	if v, found := capTree.Get(instance); found {
		return v.(capabilityMap), nil
	}

	// Find all patterns in the capability tree
	// matching the queried instance.
	merged := make(capabilityMap)
	capTree.Walk(func(pattern string, raw interface{}) bool {
		if matched, _ := path.Match(pattern, instance); matched {
			for cap := range raw.(capabilityMap) {
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
	s += fmt.Sprintf("- management = %v\n", a.management)
	for res, capTree := range a.capabilities {
		s += fmt.Sprintf("-- %s capabilities\n", res)
		capTree.Walk(func(k string, v interface{}) bool {
			s += fmt.Sprintf("    %s%s --> ", k, strings.Repeat(" ", 18-len(k)))
			for c := range v.(capabilityMap) {
				s += fmt.Sprintf("%s  ", c)
			}
			s += fmt.Sprintf("\n")
			return false
		})
	}
	return s
}

// capabilityMap is a map meant for making it easy/efficient to
// evaluate whether or not a given capability is enabled. If an
// entry exists in the map for a specific capability, it is enabled.
// Otherwise, it is not.
type capabilityMap map[string]struct{}

func (m capabilityMap) AddCapability(cap string) {
	m[cap] = struct{}{}
}

func (m capabilityMap) RemoveCapability(cap string) {
	delete(m, cap)
}

func (m capabilityMap) HasCapability(cap string) bool {
	_, found := m[cap]
	return found
}

func (m capabilityMap) Clear() {
	for cap := range m {
		delete(m, cap)
	}
}
