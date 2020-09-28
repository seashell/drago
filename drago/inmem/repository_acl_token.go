package inmem

import (
	"context"
	"strings"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	resourceTypeToken = "token"
)

// ACLTokenRepositoryAdapter :
type ACLTokenRepositoryAdapter struct {
	backend *Backend
}

// NewACLTokenRepositoryAdapter :
func NewACLTokenRepositoryAdapter(backend *Backend) domain.ACLTokenRepository {
	return &ACLTokenRepositoryAdapter{
		backend: backend,
	}
}

// GetByID ...
func (a *ACLTokenRepositoryAdapter) GetByID(ctx context.Context, id string) (*domain.ACLToken, error) {
	key := resourceKey(resourceTypeToken, id)
	if v, found := a.backend.kv[key]; found {
		return v.(*domain.ACLToken), nil
	}
	return nil, domain.ErrNotFound
}

// FindBySecret :
func (a *ACLTokenRepositoryAdapter) FindBySecret(ctx context.Context, secret string) (*domain.ACLToken, error) {
	prefix := resourcePrefix(resourceTypeToken)
	for k, v := range a.backend.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*domain.ACLToken); ok {
				if t.Secret == secret {
					return t, nil
				}
			}
		}
	}
	return nil, nil
}

// Create :
func (a *ACLTokenRepositoryAdapter) Create(ctx context.Context, t *domain.ACLToken) (*string, error) {

	id := uuid.Generate()

	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	key := resourceKey(resourceTypeToken, id)

	a.backend.kv[key] = t

	return &id, nil
}

// Update :
func (a *ACLTokenRepositoryAdapter) Update(ctx context.Context, t *domain.ACLToken) (*string, error) {

	id := t.ID
	key := resourceKey(resourceTypeToken, id)

	t.UpdatedAt = time.Now()

	a.backend.kv[key] = t

	return &id, nil
}

// DeleteByID :
func (a *ACLTokenRepositoryAdapter) DeleteByID(ctx context.Context, id string) (*string, error) {
	key := resourceKey(resourceTypeToken, id)
	delete(a.backend.kv, key)
	return &id, nil
}

// FindAll :
func (a *ACLTokenRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.ACLToken, error) {

	prefix := resourcePrefix(resourceTypeToken)

	items := []*domain.ACLToken{}
	for k, v := range a.backend.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*domain.ACLToken); ok {
				items = append(items, t)
			}
		}
	}

	return items, nil
}
