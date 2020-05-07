//go:generate echo "==> Building Web UI..."
//go:generate yarn
//go:generate yarn --cwd . build
//go:generate echo "==> Done."

//go:generate echo "==> Bundling web UI"
//go:generate go run github.com/rakyll/statik -f -src=./build
//go:generate echo "==> Done"

package ui

import (
	"net/http"

	"github.com/rakyll/statik/fs"
	_ "github.com/seashell/drago/ui/statik"
)

var Bundle http.FileSystem

func init() {
	var err error
	Bundle, err = fs.New()
	if err != nil {
		panic(err)
	}
}
