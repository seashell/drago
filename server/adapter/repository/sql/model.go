package sql

import (
	"time"
)

// Network :
type Network struct {
	ID             *string    `db:"id"`
	Name           *string    `db:"name"`
	IPAddressRange *string    `db:"ip_address_range"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

// Host :
type Host struct {
	ID               *string    `db:"id"`
	NetworkID        *string    `db:"network_id"`
	Name             *string    `db:"name"`
	IPAddress        *string    `db:"ip_address"`
	AdvertiseAddress *string    `db:"advertise_address"`
	ListenPort       *string    `db:"listen_port"`
	PublicKey        *string    `db:"public_key"`
	Table            *string    `db:"table"`
	DNS              *string    `db:"dns"`
	MTU              *string    `db:"mtu"`
	PreUp            *string    `db:"pre_up"`
	PostUp           *string    `db:"post_up"`
	PreDown          *string    `db:"pre_down"`
	PostDown         *string    `db:"post_down"`
	CreatedAt        *time.Time `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}

// Link :
type Link struct {
	ID                  *string    `db:"id"`
	NetworkID           *string    `db:"network_id"`
	FromHostID          *string    `db:"from_host_id"`
	ToHostID            *string    `db:"to_host_id"`
	AllowedIPs          []string   `db:"allowed_ips"`
	PersistentKeepalive *int       `db:"persistent_keepalive"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
}
