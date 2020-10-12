package drago

import (
	"fmt"

	http "github.com/seashell/drago/drago/infrastructure/delivery/http"
	handler "github.com/seashell/drago/drago/infrastructure/delivery/http/handler"
	middleware "github.com/seashell/drago/drago/infrastructure/delivery/http/middleware"
)

func (s *Server) setupHTTPServer() error {

	config := &http.Config{
		Logger:      s.config.Logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.HTTP),
		Handlers: map[string]http.Handler{
			"/api/acl/policies/": handler.NewACLPolicyHandlerAdapter(s.services.ACLPolicies),
			"/api/acl/tokens/":   handler.NewACLTokenHandlerAdapter(s.services.ACLTokens),
			"/api/acl/":          handler.NewACLHandlerAdapter(s.services.ACL),
			"/status":            handler.NewStatusHandlerAdapter(),
		},
		Middleware: []http.Middleware{
			middleware.CORS(),
			middleware.Logging(s.config.Logger),
		},
	}

	if s.config.UI {
		//config.Handlers["/ui/"] = handler.NewFilesystemHandlerAdapter(ui.Bundle)
		config.Handlers["/"] = handler.NewFallthroughHandlerAdapter("/ui/")
	}

	httpServer, err := http.NewServer(config)
	if err != nil {
		return err
	}

	s.httpServer = httpServer

	return nil
}
