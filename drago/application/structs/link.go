package structs

import "time"

type GetLinkInput struct {
	ID string `validate:"required,uuid4"`
	*dto
}

type GetLinkOutput struct {
	ID                  *string    `json:"id"`
	FromInterfaceID     *string    `json:"fromInterfaceId,omitempty"`
	ToInterfaceID       *string    `json:"toInterfaceId,omitempty"`
	AllowedIPs          []string   `json:"allowedIps"`
	PersistentKeepalive *int       `json:"persistentKeepalive,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `json:"updatedAt,omitempty"`
	*dto
}

type CreateLinkInput struct {
	FromInterfaceID     *string  `json:"fromInterfaceId" validate:"required,uuid4"`
	ToInterfaceID       *string  `json:"toInterfaceId" validate:"required,uuid4"`
	AllowedIPs          []string `json:"allowedIPs" validate:"dive,omitempty,cidr"`
	PersistentKeepalive *int     `json:"persistentKeepalive" validate:"omitempty,numeric"`
	*dto
}

type CreateLinkOutput struct {
	ID *string `json:"id"`
	*dto
}

type UpdateLinkInput struct {
	ID                  *string  `json:"id" validate:"required,uuid4"`
	AllowedIPs          []string `json:"allowedIPs" validate:"dive,omitempty,cidr"`
	PersistentKeepalive *int     `json:"persistentKeepalive"`
	*dto
}

type UpdateLinkOutput struct {
	ID *string `json:"id"`
	*dto
}

type DeleteLinkInput struct {
	ID string `validate:"required,uuid4"`
	*dto
}

type DeleteLinkOutput struct {
	ID *string `json:"id"`
	*dto
}

type ListLinksInput struct {
	PageInput
	NetworkIDFilter         string `query:"networkId" validate:"omitempty,uuid4"`
	SourceHostIDFilter      string `query:"fromHostId" validate:"omitempty,uuid4"`
	SourceInterfaceIDFilter string `query:"fromInterfaceId" validate:"omitempty,uuid4"`
	*dto
}

type ListLinksOutput struct {
	PageOutput
	Items []*GetLinkOutput `json:"items"`
}
