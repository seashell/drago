package structs

import (
	"time"
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
	BaseInput
	ID     string `json:"id" validate:"uuid4"`
	Secret string `json:"secret" validate:"uuid4"`
}

// ACLTokenGetOutput :
type ACLTokenGetOutput struct {
	BaseOutput
	ACLToken
}

// ACLTokenCreateInput :
type ACLTokenCreateInput struct {
	BaseInput
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type"`
	Policies []string `json:"policies"`
}

// ACLTokenCreateOutput :
type ACLTokenCreateOutput struct {
	BaseOutput
	ACLToken
}

// ACLTokenDeleteInput :
type ACLTokenDeleteInput struct {
	BaseInput
	ID string `json:"id" validate:"uuid4"`
}

// ACLTokenDeleteOutput :
type ACLTokenDeleteOutput struct {
	BaseOutput
}

// ACLTokenListInput :
type ACLTokenListInput struct {
	BaseInput
}

// ACLTokenListOutput :
type ACLTokenListOutput struct {
	BaseOutput
	Items []*ACLTokenListItem `json:"items"`
}
