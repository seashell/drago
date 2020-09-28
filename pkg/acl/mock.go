package acl

import "fmt"

type repository struct {
	tokens   []*token
	policies []*policy
}

func (r *repository) FindTokenBySecret(s string) (Token, error) {
	for _, t := range r.tokens {
		if t.secret == s {
			return t, nil
		}
	}
	return nil, nil
}

func (r *repository) GetPolicyByName(n string) (Policy, error) {
	for _, p := range r.policies {
		if p.name == n {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", n)
}

// token implements the acl.Token interface
type token struct {
	secret   string
	policies []string
}

func (t *token) Policies() []string {
	return t.policies
}

// policy implements the acl.Policy interface
type policy struct {
	name  string
	rules []Rule
}

func (p *policy) Name() string {
	return p.name
}

func (p *policy) Rules() []Rule {
	return p.rules
}

// rule implements the acl.Rule interface
type rule struct {
	resource     string
	pattern      string
	capabilities []string
}

func (r *rule) Resource() string {
	return r.resource
}

func (r *rule) Pattern() string {
	return r.pattern
}

func (r *rule) Capabilities() []string {
	return r.capabilities
}
