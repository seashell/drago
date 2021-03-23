package state

import (
	"context"

	"github.com/seashell/drago/drago/structs"
)

// Transaction :
type Transaction interface {
	Commit() (interface{}, error)
}

// Repository :
type Repository interface {
	Name() string
	Transaction(ctx context.Context) Transaction

	ACLTokenRepository
	ACLPolicyRepository

	NodeRepository

	NetworkRepository
	InterfaceRepository
	ConnectionRepository

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

// NodeRepository : Node repository interface
type NodeRepository interface {
	Nodes(ctx context.Context) ([]*structs.Node, error)
	NodeByID(ctx context.Context, id string) (*structs.Node, error)
	NodeBySecretID(ctx context.Context, sid string) (*structs.Node, error)
	UpsertNode(ctx context.Context, n *structs.Node) error
	DeleteNodes(ctx context.Context, ids []string) error
}

// InterfaceRepository : Interface repository interface
type InterfaceRepository interface {
	Interfaces(ctx context.Context) ([]*structs.Interface, error)
	InterfacesByNodeID(ctx context.Context, s string) ([]*structs.Interface, error)
	InterfacesByNetworkID(ctx context.Context, s string) ([]*structs.Interface, error)
	InterfaceByID(ctx context.Context, id string) (*structs.Interface, error)
	UpsertInterface(ctx context.Context, i *structs.Interface) error
	DeleteInterfaces(ctx context.Context, ids []string) error
}

// ConnectionRepository : Connection repository interface
type ConnectionRepository interface {
	Connections(ctx context.Context) ([]*structs.Connection, error)
	ConnectionsByNetworkID(ctx context.Context, s string) ([]*structs.Connection, error)
	ConnectionsByNodeID(ctx context.Context, s string) ([]*structs.Connection, error)
	ConnectionsByInterfaceID(ctx context.Context, s string) ([]*structs.Connection, error)
	ConnectionByInterfaceIDs(ctx context.Context, a, b string) (*structs.Connection, error)
	ConnectionByID(ctx context.Context, id string) (*structs.Connection, error)
	UpsertConnection(ctx context.Context, i *structs.Connection) error
	DeleteConnections(ctx context.Context, ids []string) error
}
