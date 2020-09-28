package domain

import (
	"context"
)

// ACLState ...
type ACLState struct {
	RootTokenID         string
	RootTokenResetIndex int
}

// Update ...
func (s *ACLState) Update(rootID string) error {
	s.RootTokenID = rootID
	s.RootTokenResetIndex++
	return nil
}

// ACLStateRepository : ACL repository interface
type ACLStateRepository interface {
	Get(ctx context.Context) (*ACLState, error)
	Set(ctx context.Context, state *ACLState) error
}
