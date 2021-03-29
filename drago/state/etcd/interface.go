package etcd

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/structs"
	"go.etcd.io/etcd/clientv3"
)

// Interfaces :
func (r *StateRepository) Interfaces(ctx context.Context) ([]*structs.Interface, error) {

	prefix := resourceKey(resourceTypeInterface, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Interface{}

	for _, el := range res.Kvs {
		iface := &structs.Interface{}
		if err := decodeValue(el.Value, iface); err != nil {
			return nil, err
		}
		items = append(items, iface)
	}

	return items, nil
}

// InterfacesByNodeID :
func (r *StateRepository) InterfacesByNodeID(ctx context.Context, id string) ([]*structs.Interface, error) {

	prefix := resourceKey(resourceTypeInterface, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Interface{}

	for _, el := range res.Kvs {
		iface := &structs.Interface{}
		if err := decodeValue(el.Value, iface); err != nil {
			return nil, err
		}
		if iface.NodeID == id {
			items = append(items, iface)
		}
	}

	return items, nil
}

// InterfacesByNetworkID :
func (r *StateRepository) InterfacesByNetworkID(ctx context.Context, id string) ([]*structs.Interface, error) {

	prefix := resourceKey(resourceTypeInterface, "")

	res, err := r.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	items := []*structs.Interface{}

	for _, el := range res.Kvs {
		iface := &structs.Interface{}
		if err := decodeValue(el.Value, iface); err != nil {
			return nil, err
		}

		if iface.NetworkID == id {
			items = append(items, iface)
		}
	}

	return items, nil
}

// InterfaceByID :
func (r *StateRepository) InterfaceByID(ctx context.Context, id string) (*structs.Interface, error) {

	key := resourceKey(resourceTypeInterface, id)

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	iface := &structs.Interface{}

	if err = decodeValue(res.Kvs[0].Value, iface); err != nil {
		return nil, err
	}

	return iface, nil
}

// UpsertInterface :
func (r *StateRepository) UpsertInterface(ctx context.Context, n *structs.Interface) error {
	key := resourceKey(resourceTypeInterface, n.ID)
	if _, err := r.client.Put(ctx, key, encodeValue(n)); err != nil {
		return err
	}
	return nil
}

// DeleteInterfaces :
func (r *StateRepository) DeleteInterfaces(ctx context.Context, ids []string) error {

	for _, id := range ids {
		key := resourceKey(resourceTypeInterface, id)
		if _, err := r.client.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}
