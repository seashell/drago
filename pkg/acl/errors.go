package acl

import (
	"errors"
)

const (
	errTokenNotFound         = "token not found"
	errPolicyNotFound        = "policy not found"
	errMissingSecretResolver = "missing secret resolver"
	errMissingPolicyResolver = "missing policy resolver"
	errResolvingSecret       = "error resolving secret"
	errResolvingPolicy       = "error resolving policy"
	errInvalidResource       = "invalid resource"
	errInvalidOperation      = "invalid operation"
	errNotAuthorized         = "not authorized"
)

var (
	// ErrTokenNotFound is returned when a token is not found
	// by the resolver function.
	ErrTokenNotFound = errors.New(errTokenNotFound)

	// ErrPolicyNotFound is returned when a policy is not found
	// by the resolver function.
	ErrPolicyNotFound = errors.New(errPolicyNotFound)

	// ErrMissingSecretResolver is returned when the no SecretResolverFunc
	// is configured.
	ErrMissingSecretResolver = errors.New(errMissingSecretResolver)

	// ErrMissingPolicyResolver is returned when the no PolicyResolverFunc
	// is configured.
	ErrMissingPolicyResolver = errors.New(errMissingPolicyResolver)

	// ErrResolvingSecret is returned when an error occurs when resolving
	// a secret.
	ErrResolvingSecret = errors.New(errResolvingSecret)

	// ErrResolvingPolicy is returned when an error occurs when resolving
	// a policy.
	ErrResolvingPolicy = errors.New(errResolvingPolicy)

	// ErrNotAuthorized is returned when the ACL does not have the
	// authorization to perform the requested operation on the specified resource.
	ErrNotAuthorized = errors.New(errNotAuthorized)

	// ErrInvalidResource is returned when the resource being queried
	// is not properly configured in the ACL system.
	ErrInvalidResource = errors.New(errInvalidResource)

	// ErrInvalidOperation is returned when the operation being queried
	// is not properly configured in the ACL system.
	ErrInvalidOperation = errors.New(errInvalidOperation)
)
