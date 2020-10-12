package structs

import "time"

// Network :
type Network struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AddressRange string    `json:"addressRange"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NetworkListItem :
type NetworkListItem struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AddressRange string    `json:"addressRange"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NetworkGetInput :
type NetworkGetInput struct {
	BaseInput
	ID string `json:"id" validate:"uuid4"`
}

// NetworkGetOutput :
type NetworkGetOutput struct {
	BaseOutput
	Network
}

// NetworkCreateInput :
type NetworkCreateInput struct {
	BaseInput
	Name         string `json:"name"`
	AddressRange string `json:"addressRange"`
}

// NetworkCreateOutput :
type NetworkCreateOutput struct {
	BaseOutput
	ID string `json:"id" validate:"uuid4"`
}

// NetworkDeleteInput :
type NetworkDeleteInput struct {
	BaseInput
	ID string `json:"id" validate:"uuid4"`
}

// NetworkDeleteOutput :
type NetworkDeleteOutput struct {
	BaseOutput
}

// NetworkListInput :
type NetworkListInput struct {
	BaseInput
}

// NetworkListOutput :
type NetworkListOutput struct {
	BaseOutput
	Items []*NetworkListItem
}
