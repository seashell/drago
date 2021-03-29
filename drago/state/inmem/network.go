package inmem

import (
	"context"
	"errors"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
)

const (
	resourceTypeNetwork = "network"
)

// Networks :
func (r *StateRepository) Networks(ctx context.Context) ([]*structs.Network, error) {
	prefix := resourcePrefix(resourceTypeNetwork)
	items := []*structs.Network{}
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if t, ok := el.Value.(*structs.Network); ok {
				items = append(items, t)
			}
		}
	}
	return items, nil
}

// NetworkByID ...
func (r *StateRepository) NetworkByID(ctx context.Context, id string) (*structs.Network, error) {
	key := resourceKey(resourceTypeNetwork, id)
	if v, found := r.kv.Get(key); found {
		return v.(*structs.Network), nil
	}
	return nil, errors.New("not found")
}

// NetworkByName ...
func (r *StateRepository) NetworkByName(ctx context.Context, name string) (*structs.Network, error) {
	prefix := resourcePrefix(resourceTypeNetwork)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if n, ok := el.Value.(*structs.Network); ok {
				if n.Name == name {
					return n, nil
				}
			}
		}
	}
	return nil, errors.New("not found")
}

// UpsertNetwork :
func (r *StateRepository) UpsertNetwork(ctx context.Context, n *structs.Network) error {
	key := resourceKey(resourceTypeNetwork, n.ID)
	r.kv.Set(key, n)
	return nil
}

// DeleteNetworks ...
func (r *StateRepository) DeleteNetworks(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeNetwork, id)
		r.kv.Delete(key)
	}
	return nil
}
