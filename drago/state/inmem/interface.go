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
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*structs.Interface); ok {
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
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if iface, ok := v.(*structs.Interface); ok {
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
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if iface, ok := v.(*structs.Interface); ok {
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
	if v, found := r.kv[key]; found {
		return v.(*structs.Interface), nil
	}
	return nil, errors.New("not found")
}

// UpsertInterface :
func (r *StateRepository) UpsertInterface(ctx context.Context, n *structs.Interface) error {
	key := resourceKey(resourceTypeInterface, n.ID)
	r.kv[key] = n
	return nil
}

// DeleteInterfaces ...
func (r *StateRepository) DeleteInterfaces(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeInterface, id)
		delete(r.kv, key)
	}
	return nil
}
