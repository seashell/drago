package api

import (
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	aclPath = "/api/acl"
)

// ACL is a handle to the ACL API
type ACL struct {
	client *Client
}

// Networks returns a handle on the networks endpoints.
func (c *Client) ACL() *ACL {
	return &ACL{client: c}
}

// Boostrap :
func (a *ACL) Bootstrap() (*structs.ACLToken, error) {

	var token structs.ACLToken
	err := a.client.createResource(path.Join(aclPath, "bootstrap"), nil, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
