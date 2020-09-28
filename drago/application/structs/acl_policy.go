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
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ACLPolicyGetInput :
type ACLPolicyGetInput struct {
	// Specifies the name of the policy.
	Name string `json:"name"`
}

// ACLPolicyGetOutput :
type ACLPolicyGetOutput struct {
	ACLPolicy
}

// ACLPolicyUpsertInput :
type ACLPolicyUpsertInput struct {
	// Specifies the name of the policy. Creates the policy if the name
	// does not exist, otherwise updates the existing policy.
	Name string `json:"name"`

	// Specifies a human readable description of the policy.
	Description string `json:"description"`

	// Specifies the Policy rules in HCL or JSON format.
	Rules string `json:"rules"`
}

// ACLPolicyUpsertOutput :
type ACLPolicyUpsertOutput struct{}

// ACLPolicyDeleteInput :
type ACLPolicyDeleteInput struct {
	Name string `json:"name"`
}

// ACLPolicyDeleteOutput :
type ACLPolicyDeleteOutput struct{}

// ACLPolicyListInput :
type ACLPolicyListInput struct{}

// ACLPolicyListOutput :
type ACLPolicyListOutput struct {
	Items []*ACLPolicyListItem `json:"items"`
}
