data_dir  = "/tmp/drago"
bind_addr = "0.0.0.0"

name = "node-2"

server {
    enabled =  false
}

client {
    enabled = true
    servers = ["192.168.100.12:8081"]
    wireguard_path = "~/wireguard"
    meta = {
        test_meta = "test_meta_value"
    }
}