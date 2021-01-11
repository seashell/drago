package etcd

import (
	"context"
	"errors"

	"github.com/seashell/drago/drago/structs"
	"go.etcd.io/etcd/clientv3"
)

// ACLPolicies :
func (r *StateRepository) ACLPolicies(ctx context.Context) ([]*structs.ACLPolicy, error) {

	prefix := resourceKey(resourceTypeACLPolicy, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.ACLPolicy{}

	for _, el := range res.Kvs {
		policy := &structs.ACLPolicy{}
		err := decodeValue(el.Value, policy)
		if err != nil {
			return nil, err
		}
		items = append(items, policy)
	}

	return items, nil
}

// ACLPolicyByName :
func (r *StateRepository) ACLPolicyByName(ctx context.Context, name string) (*structs.ACLPolicy, error) {

	key := resourceKey(resourceTypeACLPolicy, name)

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	policy := &structs.ACLPolicy{}

	err = decodeValue(res.Kvs[0].Value, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// UpsertACLPolicy :
func (r *StateRepository) UpsertACLPolicy(ctx context.Context, p *structs.ACLPolicy) error {
	key := resourceKey(resourceTypeACLPolicy, p.Name)
	_, err := r.client.Put(ctx, key, encodeValue(p))
	if err != nil {
		return err
	}
	return nil
}

// DeleteACLPolicies :
func (r *StateRepository) DeleteACLPolicies(ctx context.Context, names []string) error {
	for _, name := range names {
		key := resourceKey(resourceTypeACLPolicy, name)
		_, err := r.client.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}
