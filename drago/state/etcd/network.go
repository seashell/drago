package etcd

import (
	"context"
	"errors"
	"fmt"

	structs "github.com/seashell/drago/drago/structs"
	"go.etcd.io/etcd/clientv3"
)

// Networks :
func (r *StateRepository) Networks(ctx context.Context) ([]*structs.Network, error) {

	prefix := resourceKey(resourceTypeNetwork, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Network{}

	for _, el := range res.Kvs {
		network := &structs.Network{}
		if err := decodeValue(el.Value, network); err != nil {
			return nil, err
		}
		items = append(items, network)
	}

	return items, nil
}

// NetworkByID :
func (r *StateRepository) NetworkByID(ctx context.Context, id string) (*structs.Network, error) {

	key := resourceKey(resourceTypeNetwork, id)

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	network := &structs.Network{}
	if err = decodeValue(res.Kvs[0].Value, network); err != nil {
		return nil, err
	}

	return network, nil
}

// NetworkByName :
func (r *StateRepository) NetworkByName(ctx context.Context, s string) (*structs.Network, error) {

	prefix := resourceKey(resourceTypeNetwork, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	for _, el := range res.Kvs {
		network := &structs.Network{}
		if err := decodeValue(el.Value, network); err != nil {
			return nil, err
		}

		if network.Name == s {
			return network, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

// UpsertNetwork :
func (r *StateRepository) UpsertNetwork(ctx context.Context, n *structs.Network) error {
	key := resourceKey(resourceTypeNetwork, n.ID)
	if _, err := r.client.Put(ctx, key, encodeValue(n)); err != nil {
		return err
	}
	return nil
}

// DeleteNetworks :
func (r *StateRepository) DeleteNetworks(ctx context.Context, ids []string) error {

	for _, id := range ids {
		key := resourceKey(resourceTypeNetwork, id)
		if _, err := r.client.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}
