package etcd

import (
	"context"
	"errors"
	"fmt"

	structs "github.com/seashell/drago/drago/structs"
	"go.etcd.io/etcd/clientv3"
)

// Nodes :
func (r *StateRepository) Nodes(ctx context.Context) ([]*structs.Node, error) {

	prefix := resourceKey(resourceTypeNode, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Node{}

	for _, el := range res.Kvs {
		node := &structs.Node{}
		if err := decodeValue(el.Value, node); err != nil {
			return nil, err
		}
		items = append(items, node)
	}

	return items, nil
}

// NodeByID :
func (r *StateRepository) NodeByID(ctx context.Context, id string) (*structs.Node, error) {

	key := resourceKey(resourceTypeNode, id)

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	node := &structs.Node{}

	if err := decodeValue(res.Kvs[0].Value, node); err != nil {
		return nil, err
	}

	return node, nil
}

// Nodes :
func (r *StateRepository) NodeBySecretID(ctx context.Context, s string) (*structs.Node, error) {

	prefix := resourceKey(resourceTypeNode, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	for _, el := range res.Kvs {
		node := &structs.Node{}
		if err := decodeValue(el.Value, node); err != nil {
			return nil, err
		}
		if node.SecretID == s {
			return node, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

// UpsertNode :
func (r *StateRepository) UpsertNode(ctx context.Context, n *structs.Node) error {
	key := resourceKey(resourceTypeNode, n.ID)

	if _, err := r.client.Put(ctx, key, encodeValue(n)); err != nil {
		return err
	}
	return nil
}

// DeleteNodes :
func (r *StateRepository) DeleteNodes(ctx context.Context, ids []string) error {

	for _, id := range ids {
		key := resourceKey(resourceTypeNode, id)
		if _, err := r.client.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}
