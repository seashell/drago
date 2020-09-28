package domain

import (
	"context"
	"time"
)

const (
	// ACLTokenTypeClient ...
	ACLTokenTypeClient = "client"
	// ACLTokenTypeManagement ...
	ACLTokenTypeManagement = "management"
)

// ACLToken :
type ACLToken struct {
	ID        string
	Type      string
	Name      string
	Secret    string
	Subject   string
	Policies  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ACLTokenRepository : ACLToken repository interface
type ACLTokenRepository interface {
	GetByID(ctx context.Context, id string) (*ACLToken, error)
	FindBySecret(ctx context.Context, id string) (*ACLToken, error)
	Create(ctx context.Context, t *ACLToken) (*string, error)
	Update(ctx context.Context, t *ACLToken) (*string, error)
	DeleteByID(ctx context.Context, id string) (*string, error)
	FindAll(ctx context.Context) ([]*ACLToken, error)
}

// Merge :
func (t *ACLToken) Merge(in *ACLToken) *ACLToken {

	result := *t

	if in.Name != "" {
		result.Name = in.Name
	}
	if in.Type != "" {
		result.Type = in.Type
	}
	if in.Secret != "" {
		result.Secret = in.Secret
	}
	if in.Subject != "" {
		result.Subject = in.Subject
	}
	if in.Policies != nil {
		result.Policies = in.Policies
	}

	return &result
}
