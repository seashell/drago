package api

import (
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	aclTokensPath = "/api/acl/tokens"
)

// ACLTokens is a handle to the ACL tokens API
type ACLTokens struct {
	client *Client
}

// ACLTokens returns a handle on the ACL tokens endpoints.
func (c *Client) ACLTokens() *ACLTokens {
	return &ACLTokens{client: c}
}

// Create :
func (t *ACLTokens) Create(token *structs.ACLToken) (*structs.ACLToken, error) {

	out := &structs.ACLToken{}

	err := t.client.createResource(aclTokensPath, token, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Delete :
func (t *ACLTokens) Delete(id string) error {

	err := t.client.deleteResource(id, aclTokensPath, nil)
	if err != nil {
		return err
	}

	return nil
}

// Update :
func (t *ACLTokens) Update(token *structs.ACLToken) (*structs.ACLToken, error) {

	out := structs.ACLToken{}

	err := t.client.updateResource(token.ID, aclTokensPath, token, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Get :
func (t *ACLTokens) Get(id string) (*structs.ACLToken, error) {

	var token *structs.ACLToken
	err := t.client.getResource(aclTokensPath, id, &token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// List :
func (t *ACLTokens) List() ([]*structs.ACLTokenListStub, error) {

	var items []*structs.ACLTokenListStub
	err := t.client.listResources(path.Join(aclTokensPath, "/"), nil, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Self :
func (t *ACLTokens) Self() (*structs.ACLToken, error) {

	var token *structs.ACLToken
	err := t.client.getResource(aclTokensPath, "self", &token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
