package etcd

import (
	"context"
	"time"

	domain "github.com/seashell/drago/drago/domain"
	"go.etcd.io/etcd/clientv3"
)

const (
	resourceTypePolicy = "policy"
)

// ACLPolicyRepositoryAdapter :
type ACLPolicyRepositoryAdapter struct {
	backend     *Backend
	secretIndex map[string]string
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

	res, err := a.backend.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, domain.ErrNotFound
	}

	policy := &domain.ACLPolicy{}

	err = decodeValue(res.Kvs[0].Value, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// Upsert :
func (a *ACLPolicyRepositoryAdapter) Upsert(ctx context.Context, p *domain.ACLPolicy) (*string, error) {

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	key := resourceKey(resourceTypeToken, p.Name)

	_, err := a.backend.client.Put(ctx, key, encodeValue(p))
	if err != nil {
		return nil, err
	}

	return &p.Name, nil
}

// DeleteByName :
func (a *ACLPolicyRepositoryAdapter) DeleteByName(ctx context.Context, name string) (*string, error) {

	key := resourceKey(resourceTypePolicy, name)

	_, err := a.backend.client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	return strToPtr(name), nil
}

// FindAll :
func (a *ACLPolicyRepositoryAdapter) FindAll(ctx context.Context) ([]*domain.ACLPolicy, error) {

	prefix := resourceKey(resourceTypePolicy, "")

	res, err := a.backend.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*domain.ACLPolicy{}

	for _, el := range res.Kvs {
		token := &domain.ACLPolicy{}
		err := decodeValue(el.Value, token)
		if err != nil {
			return nil, err
		}
		items = append(items, token)
	}

	return items, nil
}
