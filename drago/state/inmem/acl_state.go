package inmem

import (
	"context"
	"errors"
	"fmt"

	structs "github.com/seashell/drago/drago/structs"
)

// ACLState :
func (r *StateRepository) ACLState(ctx context.Context) (*structs.ACLState, error) {
	key := aclStateKey()
	if v, found := r.kv[key]; found {
		return v.(*structs.ACLState), nil
	}
	return nil, errors.New("not found")
}

// ACLSetState :
func (r *StateRepository) ACLSetState(ctx context.Context, s *structs.ACLState) error {
	key := aclStateKey()
	r.kv[key] = s
	return nil
}

func aclStateKey() string {
	return fmt.Sprintf("%s/%s/global/%s", defaultPrefix, "acl", "state")
}
