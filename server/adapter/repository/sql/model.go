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
	Name             *string    `db:"name"`
	AdvertiseAddress *string    `db:"advertise_address"`
	CreatedAt        *time.Time `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}

// Interface :
type Interface struct {
	ID               *string    `db:"id"`
	Name             *string    `db:"name"`
	HostID           *string    `db:"host_id"`
	NetworkID        *string    `db:"network_id"`
	IPAddress        *string    `db:"ip_address"`
	AdvertiseAddress *string    `db:"advertise_address"`
	ListenPort       *string    `db:"listen_port"`
	Table            *string    `db:"table"`
	DNS              *string    `db:"dns"`
	MTU              *string    `db:"mtu"`
	PreUp            *string    `db:"pre_up"`
	PostUp           *string    `db:"post_up"`
	PreDown          *string    `db:"pre_down"`
	PostDown         *string    `db:"post_down"`
	PublicKey        *string    `db:"public_key"`
	CreatedAt        *time.Time `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}

// Link :
type Link struct {
	ID                  *string    `db:"id"`
	FromInterfaceID     *string    `db:"from_interface_id"`
	ToInterfaceID       *string    `db:"to_interface_id"`
	PersistentKeepalive *int       `db:"persistent_keepalive"`
	AllowedIPs          []string   `db:"-"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
}
