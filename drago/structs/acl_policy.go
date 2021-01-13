package structs

import "time"

// ACLPolicy contains a composition of subpolicies for each resource exposed by Drago.
// It can be assigned to an ACL Token and, according to the capabilities within each
// subpolicy, gives different levels of access to these resources.
type ACLPolicy struct {
	Name        string
	Description string
	Rules       []*ACLPolicyRule
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *ACLPolicy) Validate() error {
	return nil
}

func (p *ACLPolicy) Merge(in *ACLPolicy) *ACLPolicy {

	result := *p

	if in.Name != "" {
		result.Name = in.Name
	}
	if in.Description != "" {
		result.Description = in.Description
	}
	if in.Rules != nil {
		result.Rules = in.Rules
	}

	return &result
}

// Stub :
func (p *ACLPolicy) Stub() *ACLPolicyListStub {
	return &ACLPolicyListStub{
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ACLPolicyRule ...
type ACLPolicyRule struct {
	Resource     string
	Path         string
	Capabilities []string
}

// ACLPolicyListStub :
type ACLPolicyListStub struct {
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ACLPolicySpecificRequest :
type ACLPolicySpecificRequest struct {
	// Name contains the name of the policy to be retrieved.
	Name string

	QueryOptions
}

// SingleACLPolicyResponse :
type SingleACLPolicyResponse struct {
	ACLPolicy *ACLPolicy

	Response
}

// ACLPolicyUpsertRequest :
type ACLPolicyUpsertRequest struct {
	ACLPolicy *ACLPolicy

	WriteRequest
}

// ACLPolicyDeleteRequest :
type ACLPolicyDeleteRequest struct {
	// Name contains the name of the policy to be deleted.
	Names []string

	WriteRequest
}

// ACLPolicyListRequest :
type ACLPolicyListRequest struct {
	QueryOptions
}

// ACLPolicyListResponse :
type ACLPolicyListResponse struct {
	Response

	// Items contains the policies found.
	Items []*ACLPolicyListStub
}
