package auth

// Rule implements the acl.Rule interface.
type Rule struct {
	resource     string
	path         string
	capabilities []string
}

// NewRule :
func NewRule(res, path string, caps []string) *Rule {
	return &Rule{res, path, caps}
}

// Resource returns the type of resource targeted
// by the rule.
func (r *Rule) Resource() string {
	return r.resource
}

// Path returns the pattern this rule uses to
// match targeted specific resource instances.
func (r *Rule) Path() string {
	return r.path
}

// Capabilities return a slice of capabilities enabled
// on the targeted resource instances.
func (r *Rule) Capabilities() []string {
	return r.capabilities
}
