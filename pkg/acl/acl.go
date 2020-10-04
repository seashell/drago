package acl

import (
	"fmt"
	"path"
	"sort"
	"strings"

	log "github.com/seashell/drago/pkg/log"
	radix "github.com/seashell/drago/pkg/radix"
)

func init() {
	logger = defaultLogger
	resources = map[string]*resource{}
	resolveSecret = defaultResolveSecret
	resolvePolicy = defaultResolvePolicy
	anonymousToken = defaultAnonymousToken{
		policies: []string{},
	}
}

const (
	capabilityDeny = "deny"
)

// ACL is used to convert a set of policies into a structure that
// can be efficiently evaluated to determine if an action is allowed.
type ACL struct {
	management   bool
	capabilities map[string]*radix.Tree
}

// NewResource configures a resource type within the ACL system, and
// returns a Resource interface which allows for the configuration
// of capabilities and aliases associated with this resource.
func NewResource(res string) Resource {
	if _, exists := resources[res]; exists {
		panic("duplicate resource definition " + res)
	}
	r := &resource{
		name:         res,
		capabilities: map[string]*capability{},
		aliases:      map[string]*capabilityAlias{},
	}
	resources[res] = r
	return r
}

// SetLogger sets the logger to be used by the ACL system.
func SetLogger(l log.Logger) {
	logger = l
}

// AnonymousToken sets the ACL token to be used when the secret
// is empty.
func AnonymousToken(t Token) {
	anonymousToken = t
}

// SecretResolver configures how secrets are resolved to ACL tokens.
func SecretResolver(f SecretResolverFunc) {
	resolveSecret = f
}

// PolicyResolver configures how policy names are resolved to ACL policies.
func PolicyResolver(f PolicyResolverFunc) {
	resolvePolicy = f
}

// ResolveSecret creates an ACL from a secret.
func ResolveSecret(secret string) (*ACL, error) {

	var token Token

	// Handle anonymous requests.
	if secret == "" {
		token = anonymousToken
	} else {
		tkn, err := resolveSecret(secret)
		if err != nil {
			return nil, fmt.Errorf("%v : %v", ErrResolvingSecret, err)
		}
		if tkn == nil {
			return nil, ErrTokenNotFound
		}
		token = tkn
	}

	if token.IsManagement() {
		return &ACL{management: true}, nil
	}

	// Retrieve policy definitions from the repository.
	policies := []Policy{}
	for _, p := range token.Policies() {
		policy, err := resolvePolicy(p)
		if err != nil {
			return nil, fmt.Errorf("%v : %v", ErrResolvingPolicy, err)
		}
		policies = append(policies, policy)
	}

	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Name() < policies[j].Name()
	})

	// Init ACL
	acl := &ACL{capabilities: map[string]*radix.Tree{}}
	for resource := range resources {
		acl.capabilities[resource] = radix.NewTree()
	}

	// Based on the token policies, we populate the
	// capability trees for each resource type.
	for _, policy := range policies {

		for _, res := range resources {

			rules := policy.Rules()

			for _, rule := range rules {

				var capabilities capabilityMap

				if rule.Resource() == res.name {

					// Initialize capability map for this pattern if it has not yet been
					// initialized by a previously processed rule/policy.
					if leaf, found := acl.capabilities[res.name].Get(rule.Pattern()); !found {
						capabilities = make(capabilityMap)
						acl.capabilities[res.name].Set(rule.Pattern(), capabilities)
					} else {
						capabilities = leaf.(capabilityMap)
					}

					// Ignore all capabilities in the rule if a previously processed
					// rule/policy has explicitly denied access to resources matching
					// this pattern.
					if capabilities.HasCapability(capabilityDeny) {
						break
					}

					// Expand aliases before processing capabilities.
					expanded := []string{}
					for _, cap := range rule.Capabilities() {
						if res.hasAlias(cap) {
							expanded = append(expanded, res.aliases[cap].expand()...)
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

// CheckAuthorized verifies whether the ACL is authorized to perform a specific action.
// If the ACL is not authorized, an error is returned, which provides more details.
// If an operation is not explicitly enabled in the ACL, it is forbidden by default.
func (a *ACL) CheckAuthorized(res string, inst string, op string) error {

	r, ok := resources[res]
	if !ok {
		return ErrInvalidResource
	}

	if !r.hasCapability(op) {
		return ErrInvalidOperation
	}

	// A management ACL is authorized to do anything.
	if a.management {
		return nil
	}

	capabilities, err := a.queryCapabilities(res, inst)
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
func (a *ACL) queryCapabilities(resource string, instance string) (capabilityMap, error) {

	capTree, ok := a.capabilities[resource]
	if !ok {
		return nil, ErrInvalidResource
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
