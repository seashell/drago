package structs

import (
	"time"

	"github.com/seashell/drago/drago/domain"
)

var (
	// AnonymousACLToken is used no SecretID is provided, and the
	// request is made anonymously.
	AnonymousACLToken = &ACLToken{
		ID:       "anonymous",
		Name:     "Anonymous Token",
		Type:     domain.ACLTokenTypeClient,
		Policies: []string{"anonymous"},
	}
)

// ACLToken :
type ACLToken struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Secret    string    `json:"secret"`
	Type      string    `json:"type"`
	Policies  []string  `json:"policies"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ACLTokenListItem :
type ACLTokenListItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Type      string    `json:"type"`
	Policies  []string  `json:"policies"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ACLTokenGetInput :
type ACLTokenGetInput struct {
	ID     string `json:"id" validate:"uuid4"`
	Secret string `json:"secret" validate:"uuid4"`
}

// ACLTokenGetOutput :
type ACLTokenGetOutput struct {
	ACLToken
}

// ACLTokenCreateInput :
type ACLTokenCreateInput struct {
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type"`
	Policies []string `json:"policies"`
}

// ACLTokenCreateOutput :
type ACLTokenCreateOutput struct {
	ACLToken
}

// ACLTokenDeleteInput :
type ACLTokenDeleteInput struct {
	ID string `json:"id" validate:"uuid4"`
}

// ACLTokenDeleteOutput :
type ACLTokenDeleteOutput struct{}

// ACLTokenListInput :
type ACLTokenListInput struct{}

// ACLTokenListOutput :
type ACLTokenListOutput struct {
	Items []*ACLTokenListItem `json:"items"`
}
