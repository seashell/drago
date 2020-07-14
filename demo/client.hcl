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
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MTAxMTMzOTUsImlhdCI6MTU5NDc1MzM5NSwiaWQiOiJkYWVhZWUwMC02MWY3LTQ1OTAtODhiYy1jYWE3OGM2MWFmODgiLCJsYWJlbHMiOltdLCJuYmYiOjE1OTQ3NTMzOTUsInN1YiI6ImYwMzc0MWNhLTA0OWMtNGNkYy1hMmIwLWVhMzk2MzZhMjIyMSIsInR5cGUiOiJjbGllbnQifQ.xf7tM-1wD3JKW1S1U3TS0YfHIywDmfd3tiqqOa-I5zY"
    interfaces_prefix= "dg-"
    links_persistent_keepalive= 120
}

vault {
    enabled = false
}