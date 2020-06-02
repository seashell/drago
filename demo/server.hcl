data_dir  = "/opt/drago"

bind_addr = "0.0.0.0"

ui = true

server {
    enabled =  true
    data_dir = "/opt/drago/server"
    storage "inmem" {
        path = "./drago.db"
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