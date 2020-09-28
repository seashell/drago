package etcd

import (
	"context"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
	"go.etcd.io/etcd/clientv3"
)

const (
	resourceTypeToken = "token"
)

// ACLTokenRepositoryAdapter :
type ACLTokenRepositoryAdapter struct {
	backend     *Backend
	secretIndex map[string]string
}

// NewACLTokenRepositoryAdapter :
func NewACLTokenRepositoryAdapter(backend *Backend) domain.ACLTokenRepository {
	return &ACLTokenRepositoryAdapter{
		backend:     backend,
		secretIndex: map[string]string{},
	}
}

// GetByID ...
func (a *ACLTokenRepositoryAdapter) GetByID(ctx context.Context, id string) (*domain.ACLToken, error) {

	key := resourceKey(resourceTypeToken, id)

	res, err := a.backend.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, domain.ErrNotFound
	}

	token := &domain.ACLToken{}

	err = decodeValue(res.Kvs[0].Value, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// FindBySecret :
func (a *ACLTokenRepositoryAdapter) FindBySecret(ctx context.Context, secret string) (*domain.ACLToken, error) {

	prefix := resourceKey(resourceTypeToken, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, el := range res.Kvs {

		token := &domain.ACLToken{}

		err := decodeValue(el.Value, token)
		if err != nil {
			return nil, err
		}

		if token.Secret == secret {
			return token, nil
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

	_, err := a.backend.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update :
func (a *ACLTokenRepositoryAdapter) Update(ctx context.Context, t *domain.ACLToken) (*string, error) {

	id := t.ID
	key := resourceKey(resourceTypeToken, id)

	t.UpdatedAt = time.Now()

	_, err := a.backend.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// DeleteByID :
func (a *ACLTokenRepositoryAdapter) DeleteByID(ctx context.Context, id string) (*string, error) {

	key := resourceKey(resourceTypeToken, id)

	_, err := a.backend.client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	return strToPtr(id), nil
}

// FindAll :
func (a *ACLTokenRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.ACLToken, error) {

	prefix := resourceKey(resourceTypeToken, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*domain.ACLToken{}

	for _, el := range res.Kvs {
		token := &domain.ACLToken{}
		err := decodeValue(el.Value, token)
		if err != nil {
			return nil, err
		}
		items = append(items, token)
	}

	return items, nil
}
