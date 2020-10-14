package structs

import "time"

// ACLPolicy :
type ACLPolicy struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ACLPolicyListItem :
type ACLPolicyListItem struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ACLPolicyGetInput :
type ACLPolicyGetInput struct {
	BaseInput

	// Name contains the name of the policy to be retrieved.
	Name string `json:"name"`
}

// ACLPolicyGetOutput :
type ACLPolicyGetOutput struct {
	BaseOutput
	ACLPolicy
}

// ACLPolicyUpsertInput :
type ACLPolicyUpsertInput struct {
	BaseInput

	// Name contains the name of the policy to be created. If it already exists,
	// its attributes are updated.
	Name string `json:"name"`

	// Description contains a human readable description of the policy.
	Description string `json:"description"`

	// Rules contains the policy rules in HCL or JSON format.
	Rules string `json:"rules"`
}

// ACLPolicyUpsertOutput :
type ACLPolicyUpsertOutput struct {
	BaseOutput
}

// ACLPolicyDeleteInput :
type ACLPolicyDeleteInput struct {
	BaseInput

	// Name contains the name of the policy to be deleted.
	Name string `json:"name"`
}

// ACLPolicyDeleteOutput :
type ACLPolicyDeleteOutput struct {
	BaseOutput
}

// ACLPolicyListInput :
type ACLPolicyListInput struct {
	BaseInput
}

// ACLPolicyListOutput :
type ACLPolicyListOutput struct {
	BaseOutput

	// Items contains the policies found.
	Items []*ACLPolicyListItem `json:"items"`
}
