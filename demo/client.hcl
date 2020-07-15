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
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MTAyMDcwNjgsImlhdCI6MTU5NDg0NzA2OCwiaWQiOiJjMGRlN2VjOS1jNzQ0LTQzMDktYWY4Zi1jNGRhNGI0NTA1N2IiLCJsYWJlbHMiOltdLCJuYmYiOjE1OTQ4NDcwNjgsInN1YiI6IjMzMzJmNjI1LWY4M2ItNDE3Ny05Y2QwLTc3OGQ4Y2FmMmMxNSIsInR5cGUiOiJjbGllbnQifQ.X9RCk8Dl-zEweFNyBiTVWoIctrxAlyv83LesR34EwVY"
    interfaces_prefix= "dg-"
    links_persistent_keepalive= 120
}

vault {
    enabled = false
}