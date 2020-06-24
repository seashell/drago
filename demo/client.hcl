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
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MDgyMTI5ODksImlhdCI6MTU5Mjg1Mjk4OSwiaWQiOiIxZjEyODBhZS0yZWMyLTQ4MDAtOTc3OC1hZTBiZTRkMjc4NzciLCJuYmYiOjE1OTI4NTI5ODksInN1YiI6ImNmNGVlOWI5LWFmNzQtNDY0MC1iMmZkLTM2ZTJhNTZhMTYzNyIsInR5cGUiOiJjbGllbnQifQ.OX0pZxT_bSEybvLSE-VfnZbNq65150kXazxba4HYnfY"
}

vault {
    enabled = false
}