package drago

import (
	"context"
	"fmt"

	application "github.com/seashell/drago/drago/application"
	domain "github.com/seashell/drago/drago/domain"
	acl "github.com/seashell/drago/pkg/acl"
)

// ACLModel defines an ACL model containing resource types, associated
// capabilities, and aliases which can be used by the application.
func ACLModel() *acl.Model {

	model := acl.NewModel()

	model.Resource(application.ResourceACLToken).
		Capabilities(application.ACLTokenWrite, application.ACLTokenRead, application.ACLTokenList).
		Alias("read", application.ACLTokenRead, application.ACLTokenList).
		Alias("write", application.ACLTokenWrite, application.ACLTokenRead, application.ACLTokenList)

	model.Resource(application.ResourceACLPolicy).
		Capabilities(application.ACLPolicyWrite, application.ACLPolicyRead, application.ACLPolicyList).
		Alias("read", application.ACLPolicyRead, application.ACLPolicyList).
		Alias("write", application.ACLPolicyWrite, application.ACLPolicyRead, application.ACLPolicyList)

	return model
}

// AuthorizationHandlerAdapter implements the application.AuthorizationHandler
// interface, thus being capable of authorizing operations.
type AuthorizationHandlerAdapter struct {
	resolver *acl.Resolver
}

// NewAuthorizationHandlerAdapter implements the application.AuthorizationHandler
// interface and is used to enforce authorization to the application services.
func NewAuthorizationHandlerAdapter(aclTokenRepo domain.ACLTokenRepository,
	aclPolicyRepo domain.ACLPolicyRepository) *AuthorizationHandlerAdapter {

	aclResolver, _ := acl.NewResolver(&acl.ResolverConfig{
		Model: ACLModel(),
		SecretResolver: func(ctx context.Context, s string) (acl.Token, error) {
			t, err := aclTokenRepo.FindBySecret(ctx, s)
			if err != nil {
				return nil, err
			}
			return &Token{
				privileged: t.Type == domain.ACLTokenTypeManagement,
				policies:   t.Policies,
			}, nil
		},
		PolicyResolver: func(ctx context.Context, p string) (acl.Policy, error) {
			pol, err := aclPolicyRepo.GetByName(ctx, p)
			if err != nil {
				return nil, err
			}
			res := &Policy{
				name:  pol.Name,
				rules: []acl.Rule{},
			}
			for _, r := range pol.Rules {
				res.rules = append(res.rules, &Rule{
					resource:     r.Resource,
					pattern:      r.Pattern,
					capabilities: r.Capabilities,
				})
			}
			return res, nil
		},
	})

	return &AuthorizationHandlerAdapter{
		resolver: aclResolver,
	}
}

// Authorize checks whether or not the specified operation is authorized or
// not on the targeted resource and path, potentially returning an error.
func (h *AuthorizationHandlerAdapter) Authorize(ctx context.Context, sub, res, path, op string) error {
	fmt.Printf("==> op %s on %s/%s authorized!\n", op, res, path)
	return nil
}

// Policy implements the acl.Policy interface.
type Policy struct {
	name  string
	rules []acl.Rule
}

// Name returns the name of the policy.
func (p *Policy) Name() string {
	return p.name
}

// Rules return the policy rules.
func (p *Policy) Rules() []acl.Rule {
	return p.rules
}

// Rule implements the acl.Rule interface.
type Rule struct {
	resource     string
	pattern      string
	capabilities []string
}

// Resource returns the type of resource targeted
// by the rule.
func (r *Rule) Resource() string {
	return r.resource
}

// Pattern returns the pattern this rule uses to
// match targeted specific resource instances.
func (r *Rule) Pattern() string {
	return r.pattern
}

// Capabilities return a slice of capabilities enabled
// on the targeted resource instances.
func (r *Rule) Capabilities() []string {
	return r.capabilities
}

// Token implements the acl.Token interface
type Token struct {
	privileged bool
	policies   []string
}

// Policies returns a slice of policies associated with
// the token.
func (t *Token) Policies() []string {
	return t.policies
}

// IsPrivileged returns true if the token has privileged access,
// and false otherwise.
func (t *Token) IsPrivileged() bool {
	return t.privileged
}
