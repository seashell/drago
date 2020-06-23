data_dir  = "/opt/drago"

bind_addr = "0.0.0.0"

ui = true

server {
    enabled =  true
    data_dir = "/opt/drago/server"
    storage "postgresql" {
		postgresql_address  = "127.0.0.1"
	 	postgresql_port     = 5432
	 	postgresql_dbname = "drago"
	 	postgresql_user = "admin"
	 	postgresql_password = "admin"
	 	postgresql_sslmode  = "disable"
    }
}

client {
    enabled = false
}

vault {
    enabled = false
    address = "127.0.0.1/8200"
}