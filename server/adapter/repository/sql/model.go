package sql

import (
	"time"
)

// BaseModel :
type BaseModel struct {
	ID        *string   `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Network :
type Network struct {
	BaseModel
	Name           *string `db:"name"`
	IPAddressRange *string `db:"ip_address_range"`
}

// Host :
type Host struct {
	BaseModel
	NetworkID        *string  `db:"network_id"`
	LinkIDs          []string `db:"link_ids"`
	Name             *string  `db:"name"`
	IPAddress        *string  `db:"ip_address"`
	AdvertiseAddress *string  `db:"advertise_address"`
	ListenPort       *string  `db:"listen_port"`
	PublicKey        *string  `db:"public_key"`
	Table            *string  `db:"host_table"`
	DNS              *string  `db:"dns"`
	MTU              *string  `db:"mtu"`
	PreUp            *string  `db:"pre_up"`
	PostUp           *string  `db:"pre_down"`
	PreDown          *string  `db:"pre_down"`
	PostDown         *string  `db:"post_down"`
}

// Link :
type Link struct {
	BaseModel
	NetworkID           *string  `db:"network_id"`
	FromHostID          *string  `db:"from_host_id"`
	ToHostID            *string  `db:"to_host_id"`
	AllowedIPs          []string `db:"allowed_ips"`
	PersistentKeepalive *int     `db:"persistent_keepalive"`
}
