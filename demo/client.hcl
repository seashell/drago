data_dir  = "/opt/drago"

bind_addr = "0.0.0.0"

ui = false

server {
    enabled =  false
}

client {
    enabled = true
    data_dir = "/opt/drago/client"
    servers = [ "localhost:8080" ]
    token = ""
    interfaces_prefix= "dg-"
    links_persistent_keepalive= 120
}

vault {
    enabled = false
}