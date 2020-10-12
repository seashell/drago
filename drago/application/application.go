package application

import (
	"context"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

const (
	ResourceACLToken  = "token"
	ResourceACLPolicy = "policy"
	ResourceNetwork   = "network"
	ResourceHost      = "host"
)

// ACLService ...
type ACLService interface {
	Bootstrap(ctx context.Context, in *structs.ACLBootstrapInput) (*structs.ACLBootstrapOutput, error)
	ResolveToken(ctx context.Context, in *structs.ACLResolveTokenInput) (*structs.ACLResolveTokenOutput, error)
}

// ACLTokenService ...
type ACLTokenService interface {
	List(context.Context, *structs.ACLTokenListInput) (*structs.ACLTokenListOutput, error)
	GetBySecret(context.Context, *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error)
	GetByID(context.Context, *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error)
	Create(context.Context, *structs.ACLTokenCreateInput) (*structs.ACLTokenCreateOutput, error)
	Delete(context.Context, *structs.ACLTokenDeleteInput) (*structs.ACLTokenDeleteOutput, error)
}

// ACLPolicyService ...
type ACLPolicyService interface {
	List(context.Context, *structs.ACLPolicyListInput) (*structs.ACLPolicyListOutput, error)
	GetByName(context.Context, *structs.ACLPolicyGetInput) (*structs.ACLPolicyGetOutput, error)
	Upsert(context.Context, *structs.ACLPolicyUpsertInput) (*structs.ACLPolicyUpsertOutput, error)
	Delete(context.Context, *structs.ACLPolicyDeleteInput) (*structs.ACLPolicyDeleteOutput, error)
}

// NetworkService ...
type NetworkService interface {
	List(context.Context, *structs.NetworkListInput) (*structs.NetworkListOutput, error)
	GetByID(context.Context, *structs.NetworkGetInput) (*structs.NetworkGetOutput, error)
	Create(context.Context, *structs.NetworkCreateInput) (*structs.NetworkCreateOutput, error)
	Delete(context.Context, *structs.NetworkDeleteInput) (*structs.NetworkDeleteOutput, error)
}

// AuthorizationHandler specifies methods used by the
// application for authorizing access to specific resources.
type AuthorizationHandler interface {
	Authorize(ctx context.Context, subject, resource, path, op string) error
}

// Config ...
type Config struct {
	ACLEnabled          bool
	ACLStateRepository  domain.ACLStateRepository
	ACLTokenRepository  domain.ACLTokenRepository
	ACLPolicyRepository domain.ACLPolicyRepository
	NetworkRepository   domain.NetworkRepository
	HostRepository      domain.HostRepository
	AuthHandler         AuthorizationHandler
}

// Application ...
type Application struct {
	ACL         ACLService
	ACLTokens   ACLTokenService
	ACLPolicies ACLPolicyService
	Networks    NetworkService
}

// New ...
func New(config *Config) *Application {
	return &Application{
		ACL:         NewACLService(config),
		ACLTokens:   NewACLTokenService(config),
		ACLPolicies: NewACLPolicyService(config),
		Networks:    NewNetworkService(config),
	}
}
