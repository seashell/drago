package acl

// Config is used to configure the ACL system (TODO: should this be made a global!?)
type Config struct {
	// Resources contains the definitions of the resource one
	// would like to protect.
	Resources map[string]*Resource
	// Function that returns a token based on a secret,
	// or nil if it does not exist.
	ResolveToken func(string) (Token, error)
	// Function that returns a policy based on its name,
	// or an error if it does not exist.
	ResolvePolicy func(string) (Policy, error)
	// Log activity is a callback for logging any activity involving the
	// ACL. It can be used, e.g., for keeping track of the tokens used
	// for accessing each resource, or for maintaining a resour
	LogActivity func(string)
}

// Resource represents a resource in the context of which
// one wants to enforce authorization.
type Resource struct {
	Name         string
	Capabilities map[string]*Capability
	Aliases      map[string]*CapabilityAlias
}

func (r *Resource) hasCapability(s string) bool {
	_, found := r.Capabilities[s]
	return found
}

func (r *Resource) hasAlias(s string) bool {
	_, found := r.Aliases[s]
	return found
}

// Capability represents an operation possible in
// the context of a resource.
type Capability struct {
	Name        string
	Description string
}

// CapabilityAlias allows defining shorthands for
// representing multiple capabilities.
type CapabilityAlias struct {
	Label     string
	ExpandsTo []string
}

func (a *CapabilityAlias) expand() []string {
	return a.ExpandsTo
}
