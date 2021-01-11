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
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*structs.Network); ok {
				items = append(items, t)
			}
		}
	}
	return items, nil
}

// NetworkByID ...
func (r *StateRepository) NetworkByID(ctx context.Context, id string) (*structs.Network, error) {
	key := resourceKey(resourceTypeNetwork, id)
	if v, found := r.kv[key]; found {
		return v.(*structs.Network), nil
	}
	return nil, errors.New("not found")
}

// NetworkByName ...
func (r *StateRepository) NetworkByName(ctx context.Context, name string) (*structs.Network, error) {
	prefix := resourcePrefix(resourceTypeNetwork)
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if n, ok := v.(*structs.Network); ok {
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
	r.kv[key] = n
	return nil
}

// DeleteNetworks ...
func (r *StateRepository) DeleteNetworks(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeNetwork, id)
		delete(r.kv, key)
	}
	return nil
}
