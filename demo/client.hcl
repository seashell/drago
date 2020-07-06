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
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MDkyNzExMjQsImlhdCI6MTU5MzkxMTEyNCwiaWQiOiIzNzJjY2ZhYi1mZjE0LTQ0YzAtYmE3Mi00ZDQwNzI5OGEwYTYiLCJsYWJlbHMiOltdLCJuYmYiOjE1OTM5MTExMjQsInN1YiI6Ijg0NGQ4ZDg5LTM0MzAtNDJiMS1iNjdmLTYxMWY1NzdlZjMzOSIsInR5cGUiOiJjbGllbnQifQ.5W-Ct7E2VjmSFXfryV8X_5DCAbhCAKF7Rn6bgkFDSEI"
    interfaces_prefix= "dg-"
}

vault {
    enabled = false
}