package structs

import "time"

type CreateInterfaceInput struct {
	Name       *string `json:"name" validate:"required,min=1,max=50"`
	HostID     *string `json:"hostId" validate:"required,uuid4"`
	NetworkID  *string `json:"networkId" validate:"omitempty,uuid4"`
	IPAddress  *string `json:"ipAddress" validate:"omitempty,cidr"`
	ListenPort *string `json:"listenPort" validate:"omitempty,numeric,min=1,max=5"`
	PublicKey  *string `json:"publicKey" validate:""`
	Table      *string `json:"table" validate:""`
	DNS        *string `json:"dns" validate:""`
	MTU        *string `json:"mtu" validate:""`
	PreUp      *string `json:"preUp" validate:""`
	PostUp     *string `json:"postUp" validate:""`
	PreDown    *string `json:"preDown" validate:""`
	PostDown   *string `json:"postDown" validate:""`
}

type CreateInterfaceOutput struct {
	ID *string `json:"id" validate:"required,uuid4"`
}

type UpdateInterfaceInput struct {
	ID         *string `json:"id" validate:"required,uuid4"`
	NetworkID  *string `json:"networkId" validate:"omitempty,uuid4"`
	Name       *string `json:"name" validate:"min=1,max=50"`
	IPAddress  *string `json:"ipAddress" validate:"omitempty,cidr"`
	ListenPort *string `json:"listenPort" validate:"omitempty,numeric,min=1,max=5"`
	PublicKey  *string `json:"publicKey" validate:""`
	Table      *string `json:"table" validate:""`
	DNS        *string `json:"dns" validate:""`
	MTU        *string `json:"mtu" validate:""`
	PreUp      *string `json:"preUp" validate:""`
	PostUp     *string `json:"postUp" validate:""`
	PreDown    *string `json:"preDown" validate:""`
	PostDown   *string `json:"postDown" validate:""`
}

type UpdateInterfaceOutput struct {
	ID *string `json:"id" validate:"required,uuid4"`
}

type GetInterfaceInput struct {
	ID *string `json:"id" validate:"required,uuid4"`
}

type GetInterfaceOutput struct {
	ID         *string    `json:"id" validate:"required,uuid4"`
	NetworkID  *string    `json:"networkId" validate:"omitempty,uuid4"`
	Name       *string    `json:"name" validate:"min=1,max=50"`
	IPAddress  *string    `json:"ipAddress" validate:"omitempty,cidr"`
	ListenPort *string    `json:"listenPort" validate:"omitempty,numeric,min=1,max=5"`
	PublicKey  *string    `json:"publicKey" validate:""`
	Table      *string    `json:"table" validate:""`
	DNS        *string    `json:"dns" validate:""`
	MTU        *string    `json:"mtu" validate:""`
	PreUp      *string    `json:"preUp" validate:""`
	PostUp     *string    `json:"postUp" validate:""`
	PreDown    *string    `json:"preDown" validate:""`
	PostDown   *string    `json:"postDown" validate:""`
	CreatedAt  *time.Time `json:"createdAt,omitempty" validate:""`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty" validate:""`
}

type DeleteInterfaceInput struct {
	ID string `validate:"required,uuid4"`
}

type DeleteInterfaceOutput struct {
	ID string `validate:"required,uuid4"`
}

type ListInterfacesInput struct {
	PageInput
	HostIDFilter    string `query:"hostId" validate:"omitempty,uuid4"`
	NetworkIDFilter string `query:"networkId" validate:"omitempty,uuid4"`
}

type ListInterfacesOutput struct {
	PageOutput
	Items []*GetInterfaceOutput `json:"items"`
}
