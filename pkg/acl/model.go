package acl

// Model ...
type Model struct {
	resources map[string]*resource
}

// NewModel creates a new ACL model.
func NewModel() *Model {
	return &Model{
		resources: map[string]*resource{},
	}
}

// Resource configures a resource type within the ACL system, and
// returns a Resource interface which allows for the configuration
// of capabilities and aliases associated with this resource.
func (m *Model) Resource(res string) Resource {
	if _, exists := m.resources[res]; exists {
		panic("duplicate resource definition " + res)
	}
	m.resources[res] = &resource{
		name:         res,
		capabilities: map[string]*capability{},
		aliases:      map[string]*capabilityAlias{},
	}
	return m.resources[res]
}

// Resource provides functions for the configuration
// of capabilities and aliases associated to a resource.
type Resource interface {
	Capabilities(...string) Resource
	Alias(string, ...string) Resource
}

type resource struct {
	name         string
	capabilities map[string]*capability
	aliases      map[string]*capabilityAlias
}

// Capabilities ...
func (r *resource) Capabilities(caps ...string) Resource {
	for _, cap := range caps {
		r.capabilities[cap] = &capability{name: cap}
	}
	return r
}

// Alias ...
func (r *resource) Alias(alias string, caps ...string) Resource {
	// Make sure all capabilities are defined.
	for _, c := range caps {
		if !r.hasCapability(c) {
			panic("undefined capability " + c)
		}
	}
	r.aliases[alias] = &capabilityAlias{
		label:     alias,
		expandsTo: caps,
	}
	return r
}

func (r *resource) hasCapability(s string) bool {
	_, found := r.capabilities[s]
	return found
}

func (r *resource) hasAlias(s string) bool {
	_, found := r.aliases[s]
	return found
}

// capability represents an operation possible in
// the context of a resource.
type capability struct {
	name        string
	description string
}

// capabilityAlias allows defining shorthands for
// representing multiple capabilities.
type capabilityAlias struct {
	label     string
	expandsTo []string
}

func (a *capabilityAlias) expand() []string {
	return a.expandsTo
}
