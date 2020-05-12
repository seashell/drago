package static

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler :
type Handler struct {
	fs        http.FileSystem
	fsHandler http.Handler
}

// NewHandler : Create a new static files handler
func NewHandler(fs http.FileSystem) (*Handler, error) {
	return &Handler{
		fs:        fs,
		fsHandler: http.FileServer(fs),
	}, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/"
	h.fsHandler.ServeHTTP(rw, r)
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", echo.WrapHandler(h.fsHandler))
	e.GET("/static/*", echo.WrapHandler(h.fsHandler))
	e.GET("/logo.svg", echo.WrapHandler(h.fsHandler))
	e.GET("/manifest.json", echo.WrapHandler(h.fsHandler))
	e.GET("/service-worker.js", echo.WrapHandler(h.fsHandler))
	e.GET("/favicon.ico", echo.WrapHandler(h.fsHandler))
	e.GET("/*", echo.WrapHandler(h))
}
