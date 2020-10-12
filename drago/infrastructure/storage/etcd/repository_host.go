package etcd

import (
	"context"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
	"go.etcd.io/etcd/clientv3"
)

const (
	resourceTypeHost = "host"
)

// HostRepositoryAdapter :
type HostRepositoryAdapter struct {
	backend     *Backend
	secretIndex map[string]string
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

	res, err := a.backend.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, domain.ErrNotFound
	}

	host := &domain.Host{}

	err = decodeValue(res.Kvs[0].Value, host)
	if err != nil {
		return nil, err
	}

	return host, nil
}

// FindByLabels returns all hosts containing one of the labels present
// in the query.
func (a *HostRepositoryAdapter) FindByLabels(ctx context.Context, labels []string) ([]*domain.Host, error) {

	prefix := resourceKey(resourceTypeHost, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	matches := []*domain.Host{}

	for _, el := range res.Kvs {

		host := &domain.Host{}
		err := decodeValue(el.Value, host)
		if err != nil {
			return nil, err
		}

		for _, l := range labels {
			for _, hl := range host.Labels {
				if l == hl {
					matches = append(matches, host)
					break
				}
			}
		}
	}

	return matches, nil
}

// Create :
func (a *HostRepositoryAdapter) Create(ctx context.Context, t *domain.Host) (*string, error) {

	id := uuid.Generate()

	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	key := resourceKey(resourceTypeHost, id)

	_, err := a.backend.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update :
func (a *HostRepositoryAdapter) Update(ctx context.Context, t *domain.Host) (*string, error) {

	id := t.ID
	key := resourceKey(resourceTypeHost, id)

	t.UpdatedAt = time.Now()

	_, err := a.backend.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// DeleteByID :
func (a *HostRepositoryAdapter) DeleteByID(ctx context.Context, id string) (*string, error) {

	key := resourceKey(resourceTypeHost, id)

	_, err := a.backend.client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	return strToPtr(id), nil
}

// FindAll :
func (a *HostRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.Host, error) {

	prefix := resourceKey(resourceTypeHost, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*domain.Host{}

	for _, el := range res.Kvs {
		host := &domain.Host{}
		err := decodeValue(el.Value, host)
		if err != nil {
			return nil, err
		}
		items = append(items, host)
	}

	return items, nil
}
