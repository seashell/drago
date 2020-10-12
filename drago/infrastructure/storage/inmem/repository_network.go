package inmem

import (
	"context"
	"strings"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	resourceTypeNetwork = "network"
)

// NetworkRepositoryAdapter :
type NetworkRepositoryAdapter struct {
	backend *Backend
}

// NewNetworkRepositoryAdapter :
func NewNetworkRepositoryAdapter(backend *Backend) domain.NetworkRepository {
	return &NetworkRepositoryAdapter{
		backend: backend,
	}
}

// GetByID ...
func (a *NetworkRepositoryAdapter) GetByID(ctx context.Context, id string) (*domain.Network, error) {
	key := resourceKey(resourceTypeNetwork, id)
	if v, found := a.backend.kv[key]; found {
		return v.(*domain.Network), nil
	}
	return nil, domain.ErrNotFound
}

// Create :
func (a *NetworkRepositoryAdapter) Create(ctx context.Context, t *domain.Network) (*string, error) {

	id := uuid.Generate()

	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	key := resourceKey(resourceTypeNetwork, id)

	a.backend.kv[key] = t

	return &id, nil
}

// Update :
func (a *NetworkRepositoryAdapter) Update(ctx context.Context, t *domain.Network) (*string, error) {

	id := t.ID
	key := resourceKey(resourceTypeNetwork, id)

	t.UpdatedAt = time.Now()

	a.backend.kv[key] = t

	return &id, nil
}

// DeleteByID :
func (a *NetworkRepositoryAdapter) DeleteByID(ctx context.Context, id string) (*string, error) {
	key := resourceKey(resourceTypeNetwork, id)
	delete(a.backend.kv, key)
	return &id, nil
}

// FindAll :
func (a *NetworkRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.Network, error) {

	prefix := resourcePrefix(resourceTypeNetwork)

	items := []*domain.Network{}
	for k, v := range a.backend.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*domain.Network); ok {
				items = append(items, t)
			}
		}
	}

	return items, nil
}
