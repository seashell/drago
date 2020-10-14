package domain

import (
	"context"
	"time"
)

// ACLPolicy contains a composition of subpolicies for each resource exposed by Drago.
// It can be assigned to an ACL Token and, according to the capabilities within each
// subpolicy, gives different levels of access to these resources.
type ACLPolicy struct {
	Name        string
	Description string
	Rules       []*ACLPolicyRule
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ACLPolicyRule ...
type ACLPolicyRule struct {
	Resource     string
	Path         string
	Capabilities []string
}

// ACLPolicyRepository : Policy repository interface
type ACLPolicyRepository interface {
	GetByName(ctx context.Context, name string) (*ACLPolicy, error)
	Upsert(ctx context.Context, p *ACLPolicy) (*string, error)
	DeleteByName(ctx context.Context, name string) (*string, error)
	FindAll(ctx context.Context) ([]*ACLPolicy, error)
}
