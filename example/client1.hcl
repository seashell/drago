data_dir  = "/tmp/drago"
bind_addr = "0.0.0.0"

name = "node-1"

advertise { 
    peer = "192.168.100.13"
}

server {
    enabled =  false
}

client {
    enabled = true
    servers = ["192.168.99.1:8081"]
    wireguard_path = "./boringtun"
}
