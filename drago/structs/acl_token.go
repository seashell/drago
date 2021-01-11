package structs

import (
	"fmt"
	"time"
)

const (
	// ACLTokenTypeClient ...
	ACLTokenTypeClient = "client"

	// ACLTokenTypeManagement ...
	ACLTokenTypeManagement = "management"
)

// ACLToken :
type ACLToken struct {
	ID        string
	Type      string
	Name      string
	Secret    string
	Policies  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *ACLToken) Validate() error {

	if t.Type != ACLTokenTypeClient && t.Type != ACLTokenTypeManagement {
		return fmt.Errorf("invalid token type %s", t.Type)
	}

	if t.Type == ACLTokenTypeManagement && !(t.Policies == nil || len(t.Policies) == 0) {
		return fmt.Errorf("invalid token policies %v", t.Policies)
	}

	return nil
}

// Merge :
func (t *ACLToken) Merge(in *ACLToken) *ACLToken {

	result := *t

	if in.Name != "" {
		result.Name = in.Name
	}
	if in.Type != "" {
		result.Type = in.Type
	}
	if in.Secret != "" {
		result.Secret = in.Secret
	}
	if in.Policies != nil {
		result.Policies = in.Policies
	}

	return &result
}

// Stub :
func (t *ACLToken) Stub() *ACLTokenListStub {
	return &ACLTokenListStub{
		ID:        t.ID,
		Name:      t.Name,
		Type:      t.Type,
		Policies:  t.Policies,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ACLTokenListStub :
type ACLTokenListStub struct {
	ID        string
	Name      string
	Type      string
	Policies  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ACLTokenListRequest :
type ACLTokenListRequest struct {
	QueryOptions
}

// ACLTokenListResponse :
type ACLTokenListResponse struct {
	Response
	Items []*ACLTokenListStub
}

// ACLTokenSpecificRequest :
type ACLTokenSpecificRequest struct {
	QueryOptions
	ID string
}

// SingleACLTokenResponse :
type SingleACLTokenResponse struct {
	Response
	ACLToken *ACLToken
}

// ACLTokenUpsertRequest :
type ACLTokenUpsertRequest struct {
	WriteRequest
	ACLToken *ACLToken
}

// ACLTokenUpsertResponse :
type ACLTokenUpsertResponse struct {
	Response
	ACLToken *ACLToken
}

// ACLTokenDeleteRequest :
type ACLTokenDeleteRequest struct {
	WriteRequest
	IDs []string
}
