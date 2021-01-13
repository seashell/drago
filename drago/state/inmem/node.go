package inmem

import (
	"context"
	"errors"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
)

const (
	resourceTypeNode = "node"
)

// Nodes :
func (r *StateRepository) Nodes(ctx context.Context) ([]*structs.Node, error) {
	prefix := resourcePrefix(resourceTypeNode)
	items := []*structs.Node{}
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if t, ok := v.(*structs.Node); ok {
				items = append(items, t)
			}
		}
	}
	return items, nil
}

// NodeByID ...
func (r *StateRepository) NodeByID(ctx context.Context, id string) (*structs.Node, error) {
	key := resourceKey(resourceTypeNode, id)
	if v, found := r.kv[key]; found {
		return v.(*structs.Node), nil
	}
	return nil, errors.New("not found")
}

// NodeBySecretID ...
func (r *StateRepository) NodeBySecretID(ctx context.Context, s string) (*structs.Node, error) {
	prefix := resourcePrefix(resourceTypeNode)
	for k, v := range r.kv {
		if strings.HasPrefix(k, prefix) {
			if n, ok := v.(*structs.Node); ok {
				if n.SecretID == s {
					return n, nil
				}
			}
		}
	}
	return nil, errors.New("not found")
}

// UpsertNode :
func (r *StateRepository) UpsertNode(ctx context.Context, n *structs.Node) error {
	key := resourceKey(resourceTypeNode, n.ID)
	r.kv[key] = n
	return nil
}

// DeleteNodes ...
func (r *StateRepository) DeleteNodes(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeNode, id)
		delete(r.kv, key)
	}
	return nil
}
