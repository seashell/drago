package postgresql

var Migrations []string

const schema = `
CREATE TABLE IF NOT EXISTS network (
    id uuid PRIMARY KEY,
    name varchar(50) UNIQUE NOT NULL,
    ip_address_range text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS host (
    id uuid PRIMARY KEY,
    name varchar(50) NOT NULL,
    advertise_address text,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS interface (
    id uuid PRIMARY KEY,
    name varchar(32) NOT NULL,
    host_id uuid NOT NULL REFERENCES host ON DELETE CASCADE,
    network_id uuid REFERENCES network ON DELETE SET NULL,
    ip_address text,
    listen_port varchar(5),
    public_key text,
    "table" text,
    dns text,
    mtu text,
    pre_up text,
    post_up text,
    pre_down text,
    post_down text,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    UNIQUE(host_id, name)
);

CREATE TABLE IF NOT EXISTS link (
    id uuid PRIMARY KEY,
    from_interface_id uuid NOT NULL REFERENCES interface ON DELETE CASCADE,
    to_interface_id uuid NOT NULL REFERENCES interface ON DELETE CASCADE,
    allowed_ips text DEFAULT '{}',
    persistent_keepalive integer,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    UNIQUE(from_interface_id, to_interface_id)
);`

func init() {
	Migrations = []string{schema}
}
