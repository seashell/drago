package inmem

import (
	"context"
	"errors"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
)

const (
	resourceTypeInterface = "interface"
)

// Interfaces :
func (r *StateRepository) Interfaces(ctx context.Context) ([]*structs.Interface, error) {
	prefix := resourcePrefix(resourceTypeInterface)
	items := []*structs.Interface{}
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if t, ok := el.Value.(*structs.Interface); ok {
				items = append(items, t)
			}
		}
	}
	return items, nil
}

// InterfacesByNodeID ...
func (r *StateRepository) InterfacesByNodeID(ctx context.Context, s string) ([]*structs.Interface, error) {

	res := []*structs.Interface{}

	prefix := resourcePrefix(resourceTypeInterface)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if iface, ok := el.Value.(*structs.Interface); ok {
				if iface.NodeID == s {
					res = append(res, iface)
				}
			}
		}
	}
	return res, nil
}

// InterfacesByNetworkID ...
func (r *StateRepository) InterfacesByNetworkID(ctx context.Context, s string) ([]*structs.Interface, error) {

	res := []*structs.Interface{}

	prefix := resourcePrefix(resourceTypeInterface)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if iface, ok := el.Value.(*structs.Interface); ok {
				if iface.NetworkID == s {
					res = append(res, iface)
				}
			}
		}
	}
	return res, nil
}

// InterfaceByID ...
func (r *StateRepository) InterfaceByID(ctx context.Context, id string) (*structs.Interface, error) {
	key := resourceKey(resourceTypeInterface, id)
	if v, found := r.kv.Get(key); found {
		return v.(*structs.Interface), nil
	}
	return nil, errors.New("not found")
}

// UpsertInterface :
func (r *StateRepository) UpsertInterface(ctx context.Context, n *structs.Interface) error {
	key := resourceKey(resourceTypeInterface, n.ID)
	r.kv.Set(key, n)
	return nil
}

// DeleteInterfaces ...
func (r *StateRepository) DeleteInterfaces(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeInterface, id)
		r.kv.Delete(key)
	}
	return nil
}
