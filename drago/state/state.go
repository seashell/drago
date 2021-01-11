package state

import (
	"context"

	"github.com/seashell/drago/drago/structs"
)

// Repository :
type Repository interface {
	Name() string

	ACLTokenRepository
	ACLPolicyRepository
	NetworkRepository

	ACLState(ctx context.Context) (*structs.ACLState, error)
	ACLSetState(ctx context.Context, state *structs.ACLState) error
}

// ACLTokenRepository : ACLToken repository interface
type ACLTokenRepository interface {
	ACLTokens(ctx context.Context) ([]*structs.ACLToken, error)
	ACLTokenByID(ctx context.Context, id string) (*structs.ACLToken, error)
	ACLTokenBySecret(ctx context.Context, id string) (*structs.ACLToken, error)
	UpsertACLToken(ctx context.Context, t *structs.ACLToken) error
	DeleteACLTokens(ctx context.Context, ids []string) error
}

// ACLPolicyRepository : Policy repository interface
type ACLPolicyRepository interface {
	ACLPolicies(ctx context.Context) ([]*structs.ACLPolicy, error)
	ACLPolicyByName(ctx context.Context, name string) (*structs.ACLPolicy, error)
	UpsertACLPolicy(ctx context.Context, p *structs.ACLPolicy) error
	DeleteACLPolicies(ctx context.Context, names []string) error
}

// NetworkRepository : Network repository interface
type NetworkRepository interface {
	Networks(ctx context.Context) ([]*structs.Network, error)
	NetworkByID(ctx context.Context, id string) (*structs.Network, error)
	NetworkByName(ctx context.Context, name string) (*structs.Network, error)
	UpsertNetwork(ctx context.Context, n *structs.Network) error
	DeleteNetworks(ctx context.Context, ids []string) error
}
