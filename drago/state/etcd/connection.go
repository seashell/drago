package etcd

import (
	"context"
	"errors"
	"fmt"

	structs "github.com/seashell/drago/drago/structs"
	"go.etcd.io/etcd/clientv3"
)

// Connections :
func (r *StateRepository) Connections(ctx context.Context) ([]*structs.Connection, error) {

	prefix := resourceKey(resourceTypeConnection, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Connection{}

	for _, el := range res.Kvs {
		conn := &structs.Connection{}
		err := decodeValue(el.Value, conn)
		if err != nil {
			return nil, err
		}
		items = append(items, conn)
	}

	return items, nil
}

// ConnectionByID :
func (r *StateRepository) ConnectionByID(ctx context.Context, id string) (*structs.Connection, error) {

	key := resourceKey(resourceTypeConnection, id)

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	network := &structs.Connection{}

	err = decodeValue(res.Kvs[0].Value, network)
	if err != nil {
		return nil, err
	}

	return network, nil
}

// ConnectionsByNetworkID :
func (r *StateRepository) ConnectionsByNetworkID(ctx context.Context, id string) ([]*structs.Connection, error) {

	prefix := resourceKey(resourceTypeConnection, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Connection{}

	for _, el := range res.Kvs {
		conn := &structs.Connection{}
		if err := decodeValue(el.Value, conn); err != nil {
			return nil, err
		}
		if conn.NetworkID == id {
			items = append(items, conn)
		}
	}

	return items, nil
}

// ConnectionsByNodeID :
func (r *StateRepository) ConnectionsByNodeID(ctx context.Context, id string) ([]*structs.Connection, error) {

	prefix := resourceKey(resourceTypeConnection, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Connection{}

	for _, el := range res.Kvs {
		conn := &structs.Connection{}
		if err := decodeValue(el.Value, conn); err != nil {
			return nil, err
		}
		if conn.NodeIDs[0] == id || conn.NodeIDs[1] == id {
			items = append(items, conn)
		}
	}

	return items, nil
}

// ConnectionsByInterfaceID :
func (r *StateRepository) ConnectionsByInterfaceID(ctx context.Context, id string) ([]*structs.Connection, error) {

	prefix := resourceKey(resourceTypeConnection, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Connection{}

	for _, el := range res.Kvs {
		conn := &structs.Connection{}
		if err := decodeValue(el.Value, conn); err != nil {
			return nil, err
		}

		if conn.ConnectsInterface(id) {
			items = append(items, conn)
		}
	}

	return items, nil
}

// ConnectionByInterfaceIDs :
func (r *StateRepository) ConnectionByInterfaceIDs(ctx context.Context, a, b string) (*structs.Connection, error) {

	prefix := resourceKey(resourceTypeConnection, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	for _, el := range res.Kvs {
		conn := &structs.Connection{}
		if err := decodeValue(el.Value, conn); err != nil {
			return nil, err
		}

		if conn.ConnectsInterfaces(a, b) {
			return conn, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

// UpsertConnection :
func (r *StateRepository) UpsertConnection(ctx context.Context, n *structs.Connection) error {
	key := resourceKey(resourceTypeConnection, n.ID)

	_, err := r.client.Put(ctx, key, encodeValue(n))
	if err != nil {
		return err
	}
	return nil
}

// DeleteConnections :
func (r *StateRepository) DeleteConnections(ctx context.Context, ids []string) error {

	for _, id := range ids {
		key := resourceKey(resourceTypeConnection, id)
		_, err := r.client.Delete(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}
