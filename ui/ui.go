//go:generate echo "==> Building Web UI..."
//go:generate yarn
//go:generate yarn --cwd . build
//go:generate echo "==> Done."

package ui
