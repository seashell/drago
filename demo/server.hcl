data_dir  = "/opt/drago"

bind_addr = "0.0.0.0"

ui = true

server {
    enabled =  true
    data_dir = "/opt/drago/server"
    storage "postgresql" {
		postgresql_address  = "127.0.0.1"
	 	postgresql_port     = 5432
	 	postgresql_dbname = "seashell"
	 	postgresql_user = "root"
	 	postgresql_password = "password"
	 	postgresql_sslmode  = "disable"
    }
}

client {
    enabled = true
    servers = ["localhost:8080"]
    data_dir="/opt/drago/client" 
}

vault {
    enabled = false
    address = "127.0.0.1/8200"
}