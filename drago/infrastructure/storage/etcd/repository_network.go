package etcd

import (
	"context"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
	"go.etcd.io/etcd/clientv3"
)

const (
	resourceTypeNetwork = "network"
)

// NetworkRepositoryAdapter :
type NetworkRepositoryAdapter struct {
	backend     *Backend
	secretIndex map[string]string
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

	res, err := a.backend.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, domain.ErrNotFound
	}

	network := &domain.Network{}

	err = decodeValue(res.Kvs[0].Value, network)
	if err != nil {
		return nil, err
	}

	return network, nil
}

// FindBySecret :
func (a *NetworkRepositoryAdapter) FindBySecret(ctx context.Context, secret string) (*domain.Network, error) {

	prefix := resourceKey(resourceTypeNetwork, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, el := range res.Kvs {

		network := &domain.Network{}

		err := decodeValue(el.Value, network)
		if err != nil {
			return nil, err
		}

		if network.Secret == secret {
			return network, nil
		}
	}

	return nil, nil
}

// Create :
func (a *NetworkRepositoryAdapter) Create(ctx context.Context, t *domain.Network) (*string, error) {

	id := uuid.Generate()

	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	key := resourceKey(resourceTypeNetwork, id)

	_, err := a.backend.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update :
func (a *NetworkRepositoryAdapter) Update(ctx context.Context, t *domain.Network) (*string, error) {

	id := t.ID
	key := resourceKey(resourceTypeNetwork, id)

	t.UpdatedAt = time.Now()

	_, err := a.backend.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// DeleteByID :
func (a *NetworkRepositoryAdapter) DeleteByID(ctx context.Context, id string) (*string, error) {

	key := resourceKey(resourceTypeNetwork, id)

	_, err := a.backend.client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	return strToPtr(id), nil
}

// FindAll :
func (a *NetworkRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.Network, error) {

	prefix := resourceKey(resourceTypeNetwork, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*domain.Network{}

	for _, el := range res.Kvs {
		network := &domain.Network{}
		err := decodeValue(el.Value, network)
		if err != nil {
			return nil, err
		}
		items = append(items, network)
	}

	return items, nil
}
