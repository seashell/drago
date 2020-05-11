package postgresql

var Migrations []string

const schema = `CREATE TABLE IF NOT EXISTS network (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name varchar(50) UNIQUE NOT NULL,
    ip_address_range text NOT NULL
);

CREATE TABLE IF NOT EXISTS host (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name varchar(50) NOT NULL,
    ip_address text,
    advertise_address text,
    listen_port varchar(5),
    public_key text,
    "table" text,
    dns text,
    mtu text,
    pre_up text,
    post_up text,
    pre_down text,
    post_down text,
    network_id uuid NOT NULL REFERENCES network ON DELETE CASCADE,
    unique (name, network_id)
);

CREATE TABLE IF NOT EXISTS link (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    allowed_ips text DEFAULT '{}',
    persistent_keepalive integer NOT NULL,
    network_id uuid NOT NULL REFERENCES network ON DELETE CASCADE,
    to_host_id uuid NOT NULL REFERENCES host ON DELETE CASCADE,
    from_host_id uuid NOT NULL REFERENCES host ON DELETE CASCADE,
    UNIQUE(to_host_id, from_host_id)
);`

func init() {
	Migrations = []string{schema}
}
