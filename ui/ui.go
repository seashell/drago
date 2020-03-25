//go:generate echo "==> Building Web UI..."
//go:generate yarn --cwd . build
//go:generate echo "==> Done."

//go:generate echo "==> Bundling web UI"
//go:generate go run github.com/rakyll/statik -f -src=./build
//go:generate echo "==> Done"

package ui
