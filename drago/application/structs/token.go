package structs

import "time"

type CreateTokenInput struct {
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type"`
	Policies []string `json:"policies"`
}

type CreateTokenOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Secret    string    `json:"secret"`
	Type      string    `json:"type"`
	Policies  []string  `json:"policies"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateTokenInput struct {
	ID       *string  `json:"id"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type"`
	Policies []string `json:"policies"`
}

type UpdateTokenOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Secret    string    `json:"secret"`
	Type      string    `json:"type"`
	Policies  []string  `json:"policies"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetTokenInput struct {
	ID *string `json:"id" validate:"required,uuid4"`
}

type GetTokenOutput struct {
	ID        *string   `json:"id" validate:"required,uuid4"`
	Name      string    `json:"name,omitempty"`
	Secret    string    `json:"secret"`
	Type      string    `json:"type"`
	Policies  []string  `json:"policies"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListTokenInput struct {
	PageInput
}
