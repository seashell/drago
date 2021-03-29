package acl

import (
	"context"
	"fmt"
	"sort"

	log "github.com/seashell/drago/pkg/log"
	radix "github.com/seashell/drago/pkg/radix"
)

// SecretResolverFunc returns the token associated with a secret
// if it exists, or nil otherwise.
type SecretResolverFunc func(context.Context, string) (Token, error)

// PolicyResolverFunc returns a policy based on its name,
// or an error if it does not exist.
type PolicyResolverFunc func(context.Context, string) (Policy, error)

// Resolver resolves ACL secrets and policies.
type Resolver struct {
	logger        log.Logger
	model         *Model
	resolveSecret SecretResolverFunc
	resolvePolicy PolicyResolverFunc
}

// NewResolver ...
func NewResolver(config *ResolverConfig) (*Resolver, error) {

	config = DefaultResolverConfig().Merge(config)

	if config.Model == nil {
		return nil, ErrMissingModel
	}

	if config.SecretResolver == nil {
		return nil, ErrMissingSecretResolver
	}

	if config.PolicyResolver == nil {
		return nil, ErrMissingPolicyResolver
	}

	return &Resolver{
		logger:        config.Logger,
		model:         config.Model,
		resolveSecret: config.SecretResolver,
		resolvePolicy: config.PolicyResolver,
	}, nil
}

// SecretResolver configures how secrets are resolved to ACL tokens.
func (r *Resolver) SecretResolver(f SecretResolverFunc) {
	r.resolveSecret = f
}

// PolicyResolver configures how policy names are resolved to ACL policies.
func (r *Resolver) PolicyResolver(f PolicyResolverFunc) {
	r.resolvePolicy = f
}

// ResolveSecret creates an ACL from a secret.
func (r *Resolver) ResolveSecret(ctx context.Context, secret string) (*ACL, error) {

	var token Token

	// Handle anonymous requests.
	tkn, err := r.resolveSecret(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("%v : %v", ErrResolvingSecret, err)
	}
	if tkn == nil {
		return nil, ErrTokenNotFound
	}
	token = tkn

	if token.IsPrivileged() {
		return &ACL{privileged: true}, nil
	}

	// Retrieve policy definitions from the repository.
	policies := []Policy{}
	for _, p := range token.Policies() {
		policy, err := r.resolvePolicy(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("%v : %v", ErrResolvingPolicy, err)
		}
		policies = append(policies, policy)
	}

	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Name() < policies[j].Name()
	})

	// Initialize ACL struct
	acl := &ACL{capabilities: map[string]*radix.Tree{}}
	for resource := range r.model.resources {
		acl.capabilities[resource] = radix.NewTree()
	}

	// Based on the token policies, we populate the
	// capability trees for each resource type.
	for _, policy := range policies {

		for _, res := range r.model.resources {

			rules := policy.Rules()

			for _, rule := range rules {

				var capabilities capMap

				if rule.Resource() == res.name {

					// Initialize capability map for this pattern if it has not yet been
					// initialized by a previously processed rule/policy.
					if leaf, found := acl.capabilities[res.name].Get(rule.Path()); !found {
						capabilities = make(capMap)
						acl.capabilities[res.name].Set(rule.Path(), capabilities)
					} else {
						capabilities = leaf.(capMap)
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
