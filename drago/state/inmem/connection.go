package inmem

import (
	"context"
	"errors"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
)

const (
	resourceTypeConnection = "connection"
)

// Connections :
func (r *StateRepository) Connections(ctx context.Context) ([]*structs.Connection, error) {
	prefix := resourcePrefix(resourceTypeConnection)
	items := []*structs.Connection{}
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if t, ok := el.Value.(*structs.Connection); ok {
				items = append(items, t)
			}
		}
	}
	return items, nil
}

// ConnectionByID ...
func (r *StateRepository) ConnectionByID(ctx context.Context, id string) (*structs.Connection, error) {
	key := resourceKey(resourceTypeConnection, id)
	if v, found := r.kv.Get(key); found {
		return v.(*structs.Connection), nil
	}
	return nil, errors.New("not found")
}

// ConnectionByInterfaceIDs ...
func (r *StateRepository) ConnectionByInterfaceIDs(ctx context.Context, a, b string) (*structs.Connection, error) {
	prefix := resourcePrefix(resourceTypeConnection)

	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if c, ok := el.Value.(*structs.Connection); ok {
				if c.ConnectsInterfaces(a, b) {
					return c, nil
				}
			}
		}
	}

	return nil, errors.New("not found")
}

// ConnectionsByNetworkID ...
func (r *StateRepository) ConnectionsByNetworkID(ctx context.Context, id string) ([]*structs.Connection, error) {

	res := []*structs.Connection{}

	prefix := resourcePrefix(resourceTypeConnection)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if conn, ok := el.Value.(*structs.Connection); ok {
				if conn.NetworkID == id {
					res = append(res, conn)
				}
			}
		}
	}

	return res, nil
}

// ConnectionsByNodeID ...
func (r *StateRepository) ConnectionsByNodeID(ctx context.Context, id string) ([]*structs.Connection, error) {

	res := []*structs.Connection{}

	prefix := resourcePrefix(resourceTypeConnection)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if conn, ok := el.Value.(*structs.Connection); ok {
				if conn.PeerSettings[0].NodeID == id || conn.PeerSettings[1].NodeID == id {
					res = append(res, conn)
				}
			}
		}
	}

	return res, nil
}

// ConnectionsByInterfaceID ...
func (r *StateRepository) ConnectionsByInterfaceID(ctx context.Context, id string) ([]*structs.Connection, error) {

	res := []*structs.Connection{}

	prefix := resourcePrefix(resourceTypeConnection)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if conn, ok := el.Value.(*structs.Connection); ok {
				if conn.ConnectsInterface(id) {
					res = append(res, conn)
				}
			}
		}
	}

	return res, nil
}

// UpsertConnection :
func (r *StateRepository) UpsertConnection(ctx context.Context, n *structs.Connection) error {
	key := resourceKey(resourceTypeConnection, n.ID)
	r.kv.Set(key, n)
	return nil
}

// DeleteConnections ...
func (r *StateRepository) DeleteConnections(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeConnection, id)
		r.kv.Delete(key)
	}
	return nil
}
