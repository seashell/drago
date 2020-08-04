package structs

import "time"

type CreateHostInput struct {
	*dto
	Name             *string  `json:"name" validate:"required,min=1,max=50"`
	AdvertiseAddress *string  `json:"advertiseAddress" validate:"omitempty,ip_addr|fqdn"`
	Labels           []string `json:"labels" validate:"dive,omitempty,dashedalphanum"`
}

type CreateHostOutput struct {
	*dto
	ID *string `json:"id"`
}

type UpdateHostInput struct {
	*dto
	ID               string   `validate:"required,uuid4"`
	Name             *string  `json:"name,omitempty"`
	AdvertiseAddress *string  `json:"advertiseAddress,omitempty"`
	Labels           []string `json:"labels,omitempty"`
}

type UpdateHostOutput struct {
	*dto
	ID *string `json:"id"`
}

type GetHostInput struct {
	*dto
	ID string `validate:"required,uuid4"`
}

type GetHostOutput struct {
	*dto
	ID               string     `validate:"required,uuid4"`
	Name             *string    `json:"name,omitempty"`
	AdvertiseAddress *string    `json:"advertiseAddress,omitempty"`
	Labels           []string   `json:"labels,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

type ListHostsInput struct {
	*dto
	PageInput
}

type ListHostsOutput struct {
	*dto
	PageOutput
	Items []*GetHostOutput `json:"items"`
}
