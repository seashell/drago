package etcd

import (
	"context"
	"errors"

	"github.com/seashell/drago/drago/structs"
	"go.etcd.io/etcd/clientv3"
)

const (
	resourceTypeToken = "token"
)

// ACLTokens :
func (r *StateRepository) ACLTokens(ctx context.Context) ([]*structs.ACLToken, error) {

	prefix := resourceKey(resourceTypeToken, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.ACLToken{}

	for _, el := range res.Kvs {
		token := &structs.ACLToken{}
		err := decodeValue(el.Value, token)
		if err != nil {
			return nil, err
		}
		items = append(items, token)
	}

	return items, nil
}

// ACLTokenByID :
func (r *StateRepository) ACLTokenByID(ctx context.Context, id string) (*structs.ACLToken, error) {

	key := resourceKey(resourceTypeToken, id)

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	token := &structs.ACLToken{}

	err = decodeValue(res.Kvs[0].Value, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ACLTokenBySecret :
func (r *StateRepository) ACLTokenBySecret(ctx context.Context, secret string) (*structs.ACLToken, error) {

	prefix := resourceKey(resourceTypeToken, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, el := range res.Kvs {

		token := &structs.ACLToken{}

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

// UpsertACLToken :
func (r *StateRepository) UpsertACLToken(ctx context.Context, t *structs.ACLToken) error {
	key := resourceKey(resourceTypeToken, t.ID)
	_, err := r.client.Put(ctx, key, encodeValue(t))
	if err != nil {
		return err
	}
	return nil
}

// DeleteACLTokens :
func (r *StateRepository) DeleteACLTokens(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeToken, id)
		_, err := r.client.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}
