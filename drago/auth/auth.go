package auth

import (
	"context"

	"github.com/seashell/drago/pkg/acl"
)

// AuthorizationHandler abstracts a handler capable of
// authorizing subjects to perform specific operations on
// resources at a given path
type AuthorizationHandler interface {
	Authorize(ctx context.Context, sub, res, path, op string) error
}

// authorizationHandler implements the AuthorizationHandler
// interface, thus being capable of authorizing operations.
type authorizationHandler struct {
	resolver *acl.Resolver
}

// NewAuthorizationHandler returns a new AuthorizationHandler
func NewAuthorizationHandler(
	model *acl.Model,
	secretResolver acl.SecretResolverFunc,
	policyResolver acl.PolicyResolverFunc) AuthorizationHandler {

	aclResolver, _ := acl.NewResolver(&acl.ResolverConfig{
		Model:          model,
		SecretResolver: secretResolver,
		PolicyResolver: policyResolver,
	})

	return &authorizationHandler{
		resolver: aclResolver,
	}
}

// Authorize checks whether or not the specified operation is authorized or
// not on the targeted resource and path, potentially returning an error.
func (h *authorizationHandler) Authorize(ctx context.Context, sub, res, path, op string) error {
	acl, err := h.resolver.ResolveSecret(ctx, sub)
	if err != nil {
		return err
	}
	return acl.CheckAuthorized(ctx, res, path, op)
}
