package etcd

import (
	"context"
	"errors"

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
		err := decodeValue(el.Value, network)
		if err != nil {
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

	err = decodeValue(res.Kvs[0].Value, network)
	if err != nil {
		return nil, err
	}

	return network, nil
}

// UpsertNetwork :
func (r *StateRepository) UpsertNetwork(ctx context.Context, n *structs.Network) error {
	key := resourceKey(resourceTypeNetwork, n.ID)

	_, err := r.client.Put(ctx, key, encodeValue(n))
	if err != nil {
		return err
	}
	return nil
}

// DeleteNetworks :
func (r *StateRepository) DeleteNetworks(ctx context.Context, ids []string) error {

	for _, id := range ids {
		key := resourceKey(resourceTypeNetwork, id)
		_, err := r.client.Delete(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}
