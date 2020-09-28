package etcd

import (
	"context"
	"fmt"

	domain "github.com/seashell/drago/drago/domain"
)

// ACLStateRepositoryAdapter :
type ACLStateRepositoryAdapter struct {
	backend *Backend
}

// NewACLStateRepositoryAdapter :
func NewACLStateRepositoryAdapter(backend *Backend) domain.ACLStateRepository {
	return &ACLStateRepositoryAdapter{
		backend: backend,
	}
}

// Get ...
func (a *ACLStateRepositoryAdapter) Get(ctx context.Context) (*domain.ACLState, error) {

	key := aclStateKey()

	res, err := a.backend.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, domain.ErrNotFound
	}

	state := &domain.ACLState{}

	err = decodeValue(res.Kvs[0].Value, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// Save :
func (a *ACLStateRepositoryAdapter) Save(ctx context.Context, s *domain.ACLState) error {

	key := aclStateKey()

	_, err := a.backend.client.Put(ctx, key, encodeValue(s))
	if err != nil {
		return err
	}

	return nil
}

func aclStateKey() string {
	return fmt.Sprintf("%s/%s/global/%s", defaultPrefix, "acl", "state")
}
