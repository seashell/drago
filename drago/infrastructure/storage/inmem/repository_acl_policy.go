package inmem

import (
	"context"
	"strings"
	"time"

	domain "github.com/seashell/drago/drago/domain"
)

const (
	resourceTypePolicy = "policy"
)

// ACLPolicyRepositoryAdapter :
type ACLPolicyRepositoryAdapter struct {
	backend *Backend
}

// NewACLPolicyRepositoryAdapter :
func NewACLPolicyRepositoryAdapter(backend *Backend) domain.ACLPolicyRepository {
	return &ACLPolicyRepositoryAdapter{
		backend: backend,
	}
}

// GetByName :
func (a *ACLPolicyRepositoryAdapter) GetByName(ctx context.Context, name string) (*domain.ACLPolicy, error) {
	key := resourceKey(resourceTypePolicy, name)
	if v, found := a.backend.kv[key]; found {
		return v.(*domain.ACLPolicy), nil
	}
	return nil, domain.ErrNotFound
}

// Upsert :
func (a *ACLPolicyRepositoryAdapter) Upsert(ctx context.Context, p *domain.ACLPolicy) (*string, error) {

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	key := resourceKey(resourceTypePolicy, p.Name)

	a.backend.kv[key] = p

	return &p.Name, nil
}

// DeleteByName :
func (a *ACLPolicyRepositoryAdapter) DeleteByName(ctx context.Context, name string) (*string, error) {
	key := resourceKey(resourceTypePolicy, name)
	delete(a.backend.kv, key)
	return &name, nil
}

// FindAll :
func (a *ACLPolicyRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.ACLPolicy, error) {
	prefix := resourcePrefix(resourceTypePolicy)

	items := []*domain.ACLPolicy{}
	for k, v := range a.backend.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*domain.ACLPolicy); ok {
				items = append(items, t)
			}
		}
	}

	return items, nil
}
