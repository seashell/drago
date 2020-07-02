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

	// Import the UI bundle
	_ "github.com/seashell/drago/ui/statik"
)

// Bundle containing pre-built SPA
var Bundle http.FileSystem

func init() {
	var err error
	Bundle, err = fs.New()
	if err != nil {
		panic(err)
	}
}
