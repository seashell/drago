package client

const tmpl = `[Interface]
Address = {{ .Interface.Address }}
ListenPort = {{ .Interface.ListenPort }}
PrivateKey = {{ .Interface.PrivateKey }}
{{with .Peers}}
{{ range . -}}
{{ if .Endpoint -}}
[Peer]
Endpoint = {{ .Endpoint }}
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllowedIPs }}
PersistentKeepalive = {{ .PersistentKeepalive }}
{{else -}}
[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllowedIPs }}
PersistentKeepalive = {{ .PersistentKeepalive }}
{{end}}
{{end}}
{{end}}
`
