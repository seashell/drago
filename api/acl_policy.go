package api

import (
	"fmt"
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	aclPoliciesPath = "/api/acl/policies"
)

// ACLPolicies is a handle to the ACL policies API
type ACLPolicies struct {
	client *Client
}

// ACLPolicies returns a handle on the ACL policies endpoints.
func (c *Client) ACLPolicies() *ACLPolicies {
	return &ACLPolicies{client: c}
}

// Create :
func (p *ACLPolicies) Upsert(policy *structs.ACLPolicy) error {
	return fmt.Errorf("not implemented")
}

// Delete :
func (p *ACLPolicies) Delete(name string) error {

	err := p.client.deleteResource(name, aclPoliciesPath, nil)
	if err != nil {
		return err
	}

	return nil
}

// Get :
func (p *ACLPolicies) Get(name string) (*structs.ACLPolicy, error) {

	var policy *structs.ACLPolicy
	err := p.client.getResource(aclPoliciesPath, name, &policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// List :
func (p *ACLPolicies) List() ([]*structs.ACLPolicyListStub, error) {

	var items []*structs.ACLPolicyListStub
	err := p.client.listResources(path.Join(aclPoliciesPath, "/"), nil, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
