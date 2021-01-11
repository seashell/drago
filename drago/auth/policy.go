package auth

import "github.com/seashell/drago/pkg/acl"

// Policy implements the acl.Policy interface.
type Policy struct {
	name  string
	rules []acl.Rule
}

// NewPolicy :
func NewPolicy(name string, rules []acl.Rule) *Policy {
	return &Policy{name, rules}
}

// Name returns the name of the policy.
func (p *Policy) Name() string {
	return p.name
}

// Rules return the policy rules.
func (p *Policy) Rules() []acl.Rule {
	return p.rules
}

// AddRule appends an acl.Rule to the policy rules, returning the resulting slice.
func (p *Policy) AddRule(r acl.Rule) []acl.Rule {
	p.rules = append(p.rules, r)
	return p.rules
}
