package dac

type Network struct {
	ID             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	IPAddressRange *string `json:"ipAddressRange,omitempty"`
}

type Host struct {
	ID               *string  `json:"id,omitempty"`
	Name             *string  `json:"name"`
	AdvertiseAddress *string  `json:"advertiseAddress,omitempty"`
	Labels           []string `json:"labels"`
}

type Interface struct {
	ID         *string `json:"id,omitempty"`
	Name       *string `json:"name"`
	HostID     *string `json:"hostId"`
	NetworkID  *string `json:"networkId,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
	ListenPort *string `json:"listenPort,omitempty"`
	PublicKey  *string `json:"publicKey"`
	Table      *string `json:"table"`
	DNS        *string `json:"dns"`
	MTU        *string `json:"mtu"`
	PreUp      *string `json:"preUp"`
	PostUp     *string `json:"postUp"`
	PreDown    *string `json:"preDown"`
	PostDown   *string `json:"postDown"`
}

type Link struct {
	ID                  *string  `json:"id,omitempty"`
	FromInterfaceID     *string  `json:"fromInterfaceId"`
	ToInterfaceID       *string  `json:"toInterfaceId"`
	AllowedIPs          []string `json:"allowedIPs"`
	PersistentKeepalive *int     `json:"persistentKeepalive,omitempty"`
}

type Token struct {
	ID        *string  `json:"id,omitempty"`
	Type      *string  `json:"type,omitempty"`
	Subject   *string  `json:"subject,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Raw       *string  `json:"secret,omitempty"`
	IssuedAt  *int64   `json:"issuedAt,omitempty"`
	ExpiresAt *int64   `json:"expiresAt,omitempty"`
	NotBefore *int64   `json:"notBefore,omitempty"`
}

type Page struct {
	Page       *int `json:"page,omitempty"`
	PerPage    *int `json:"perPage,omitempty"`
	PageCount  *int `json:"pageCount,omitempty"`
	TotalCount *int `json:"totalCount,omitempty"`
}
