data_dir  = "/tmp/drago"
bind_addr = "0.0.0.0"

name = "node-2"

advertise { 
    peer = "192.168.100.14"
}

server {
    enabled =  false
}

client {
    enabled = true
    servers = ["192.168.100.12:8081"]
    wireguard_path = "/home/eschmidt/wireguard"
    meta = {
        test_meta = "test_meta_value"
    }
}