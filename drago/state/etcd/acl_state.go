package etcd

import (
	"context"
	"errors"
	"fmt"

	"github.com/seashell/drago/drago/structs"
)

// ACLState :
func (r *StateRepository) ACLState(ctx context.Context) (*structs.ACLState, error) {

	key := aclStateKey()

	res, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return nil, errors.New("not found")
	}

	state := &structs.ACLState{}

	err = decodeValue(res.Kvs[0].Value, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// ACLSetState :
func (r *StateRepository) ACLSetState(ctx context.Context, s *structs.ACLState) error {

	key := aclStateKey()

	_, err := r.client.Put(ctx, key, encodeValue(s))
	if err != nil {
		return err
	}

	return nil
}

func aclStateKey() string {
	return fmt.Sprintf("%s/%s/global/%s", defaultPrefix, "acl", "state")
}
