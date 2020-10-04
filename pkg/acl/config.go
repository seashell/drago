package acl

import (
	"github.com/seashell/drago/pkg/log"
)

var logger log.Logger
var resources map[string]*resource
var resolveSecret SecretResolverFunc
var resolvePolicy PolicyResolverFunc
var anonymousToken Token

var defaultLogger = &simpleLogger{
	fields: map[string]interface{}{},
	options: log.LoggerOptions{
		Prefix: "acl",
		Level:  levelInfo,
	},
}

var defaultResolveSecret = func(string) (Token, error) {
	return nil, ErrMissingSecretResolver
}

var defaultResolvePolicy = func(string) (Policy, error) {
	return nil, ErrMissingPolicyResolver
}

type defaultAnonymousToken struct {
	policies []string
}

func (t defaultAnonymousToken) Policies() []string {
	return t.policies
}

func (t defaultAnonymousToken) IsManagement() bool {
	return false
}

// SecretResolverFunc returns the token associated with a secret
// if it exists, or nil otherwise.
type SecretResolverFunc func(string) (Token, error)

// PolicyResolverFunc returns a policy based on its name,
// or an error if it does not exist.
type PolicyResolverFunc func(string) (Policy, error)

// Resource exposes functions for the configuration
// of capabilities and aliases associated to a resource.
type Resource interface {
	AddCapability(...string) Resource
	AddAlias(string, ...string) Resource
}

// Resource represents a resource in the context of which
// one wants to enforce authorization.
type resource struct {
	name         string
	capabilities map[string]*capability
	aliases      map[string]*capabilityAlias
}

// Capability ...
func (r *resource) AddCapability(caps ...string) Resource {
	for _, cap := range caps {
		r.capabilities[cap] = &capability{name: cap}
	}
	return r
}

// Alias ...
func (r *resource) AddAlias(alias string, caps ...string) Resource {
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
