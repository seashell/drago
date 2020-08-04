package structs

import "time"

type CreateNetworkInput struct {
	Name           *string `json:"name,omitempty"`
	IPAddressRange *string `json:"ipAddressRange,omitempty"`
}

type CreateNetworkOutput struct {
	ID *string `json:"id"`
}

type UpdateNetworkInput struct {
	ID             *string `json:"id"`
	Name           *string `json:"name,omitempty"`
	IPAddressRange *string `json:"ipAddressRange,omitempty"`
}

type UpdateNetworkOutput struct {
	ID *string `json:"id"`
}

type GetNetworkInput struct {
	ID *string `validate:"required,uuid4"`
}

type GetNetworkOutput struct {
	ID             *string    `json:"id"`
	Name           *string    `json:"name,omitempty"`
	IPAddressRange *string    `json:"ipAddressRange,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
}

type DeleteNetworkInput struct {
	ID string `validate:"required,uuid4"`
}

type DeleteNetworkOutput struct {
	ID string `validate:"required,uuid4"`
}

type ListNetworksInput struct {
	PageInput
}

type ListNetworksOutput struct {
	*dto
	PageOutput
	Items []*GetNetworkOutput `json:"items"`
}
