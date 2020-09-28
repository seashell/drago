package inmem

import (
	"context"
	"fmt"

	domain "github.com/seashell/drago/drago/domain"
)

type ACLStateRepositoryAdapter struct {
	backend *Backend
}

func NewACLStateRepositoryAdapter(b *Backend) domain.ACLStateRepository {
	return &ACLStateRepositoryAdapter{
		backend: b,
	}
}

func (a *ACLStateRepositoryAdapter) Get(ctx context.Context) (*domain.ACLState, error) {

	key := aclStateKey()

	if v, found := a.backend.kv[key]; found {
		return v.(*domain.ACLState), nil
	}

	return nil, domain.ErrNotFound
}

func (a *ACLStateRepositoryAdapter) Set(ctx context.Context, s *domain.ACLState) error {
	key := aclStateKey()
	a.backend.kv[key] = s
	return nil
}

func aclStateKey() string {
	return fmt.Sprintf("%s/%s/global/%s", defaultPrefix, "acl", "state")
}
