package client

import (
	"bytes"
	"path/filepath"
	"html/template"
	"os"
)

const tmpl = `[Interface]
{{ if .Interface.ListenPort -}}
ListenPort = {{ .Interface.ListenPort }}
{{end -}}
PrivateKey = {{ .Keys.PrivateKey }}
{{with .Peers}}
{{ range . -}}
[Peer]
{{ if .Endpoint -}}
Endpoint = {{ .Endpoint }}
{{end -}}
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllowedIPs }}
PersistentKeepalive = {{ .PersistentKeepalive }}
{{end}}
{{end}}
`

func templateToFile(tmpl string, path string, ctx interface{}) error {
	
	dir,_ := filepath.Split(path)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	t := template.Must(template.New("").Parse(tmpl))

	err = t.Execute(f, ctx)
	if err != nil {
		return err
	}

	return nil
}

func templateToString(tmpl string, ctx interface{}) (string, error) {

	t := template.Must(template.New("").Parse(tmpl))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, ctx); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
