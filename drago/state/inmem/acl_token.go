package inmem

import (
	"context"
	"errors"
	"strings"

	structs "github.com/seashell/drago/drago/structs"
)

const (
	resourceTypeToken = "token"
)

// ACLTokenByID ...
func (r *StateRepository) ACLTokenByID(ctx context.Context, id string) (*structs.ACLToken, error) {
	key := resourceKey(resourceTypeToken, id)
	if v, found := r.kv.Get(key); found {
		return v.(*structs.ACLToken), nil
	}
	return nil, errors.New("not found")
}

// ACLTokenBySecret :
func (r *StateRepository) ACLTokenBySecret(ctx context.Context, secret string) (*structs.ACLToken, error) {
	prefix := resourcePrefix(resourceTypeToken)
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if t, ok := el.Value.(*structs.ACLToken); ok {
				if t.Secret == secret {
					return t, nil
				}
			}
		}
	}
	return nil, nil
}

// UpsertACLToken :
func (r *StateRepository) UpsertACLToken(ctx context.Context, t *structs.ACLToken) error {
	key := resourceKey(resourceTypeToken, t.ID)
	r.kv.Set(key, t)
	return nil
}

// DeleteACLTokens :
func (r *StateRepository) DeleteACLTokens(ctx context.Context, ids []string) error {
	for _, id := range ids {
		key := resourceKey(resourceTypeToken, id)
		r.kv.Delete(key)
	}
	return nil
}

// ACLTokens :
func (r *StateRepository) ACLTokens(ctx context.Context) ([]*structs.ACLToken, error) {
	prefix := resourcePrefix(resourceTypeToken)
	items := []*structs.ACLToken{}
	for el := range r.kv.Iter() {
		if strings.HasPrefix(el.Key, prefix) {
			if t, ok := el.Value.(*structs.ACLToken); ok {
				items = append(items, t)
			}
		}
	}
	return items, nil
}
