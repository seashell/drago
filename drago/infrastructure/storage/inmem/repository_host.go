package inmem

import (
	"context"
	"strings"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	resourceTypeHost = "host"
)

// HostRepositoryAdapter :
type HostRepositoryAdapter struct {
	backend *Backend
}

// NewHostRepositoryAdapter :
func NewHostRepositoryAdapter(backend *Backend) domain.HostRepository {
	return &HostRepositoryAdapter{
		backend: backend,
	}
}

// GetByID ...
func (a *HostRepositoryAdapter) GetByID(ctx context.Context, id string) (*domain.Host, error) {
	key := resourceKey(resourceTypeHost, id)
	if v, found := a.backend.kv[key]; found {
		return v.(*domain.Host), nil
	}
	return nil, domain.ErrNotFound
}

// Create :
func (a *HostRepositoryAdapter) Create(ctx context.Context, t *domain.Host) (*string, error) {

	id := uuid.Generate()

	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	key := resourceKey(resourceTypeHost, id)

	a.backend.kv[key] = t

	return &id, nil
}

// Update :
func (a *HostRepositoryAdapter) Update(ctx context.Context, t *domain.Host) (*string, error) {

	id := t.ID
	key := resourceKey(resourceTypeHost, id)

	t.UpdatedAt = time.Now()

	a.backend.kv[key] = t

	return &id, nil
}

// DeleteByID :
func (a *HostRepositoryAdapter) DeleteByID(ctx context.Context, id string) (*string, error) {
	key := resourceKey(resourceTypeHost, id)
	delete(a.backend.kv, key)
	return &id, nil
}

// FindAll :
func (a *HostRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.Host, error) {

	prefix := resourcePrefix(resourceTypeHost)

	items := []*domain.Host{}
	for k, v := range a.backend.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*domain.Host); ok {
				items = append(items, t)
			}
		}
	}

	return items, nil
}

// FindByLabels : TODO
func (a *HostRepositoryAdapter) FindByLabels(ctx context.Context, labels []string) ([]*domain.Host, error) {

	prefix := resourcePrefix(resourceTypeHost)

	items := []*domain.Host{}
	for k, v := range a.backend.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*domain.Host); ok {
				items = append(items, t)
			}
		}
	}

	return items, nil
}
