package spa

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler :
type Handler struct {
	fs        http.FileSystem
	fsHandler http.Handler
}

// NewHandler : Create a new SPA handler
func NewHandler(fs http.FileSystem) (*Handler, error) {
	return &Handler{
		fs:        fs,
		fsHandler: http.FileServer(fs),
	}, nil
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	ui := e.Group("/ui/")

	ui.GET("", echo.WrapHandler(http.StripPrefix("/ui/", h.fsHandler)))
	ui.GET("static/*", echo.WrapHandler(http.StripPrefix("/ui/", h.fsHandler)))

	ui.GET("*", echo.WrapHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		h.fsHandler.ServeHTTP(rw, r)
	})))

	// Root fallthrough
	e.GET("/", echo.WrapHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(rw, r, "/ui/", http.StatusTemporaryRedirect)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	})))

}
