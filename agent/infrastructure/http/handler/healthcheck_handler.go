package handler

import (
	http "github.com/seashell/drago/agent/infrastructure/http"
	log "github.com/seashell/drago/pkg/log"
)

type HealthcheckHandlerAdapter struct {
	http.BaseHandlerAdapter
	logger log.Logger
}

func NewHealthcheckHandlerAdapter(logger log.Logger) *NetworkHandlerAdapter {

	a := &NetworkHandlerAdapter{}
	a.logger = logger

	a.RegisterHandlerFunc("GET", "/", a.healthcheck)

	return a
}

func (a *NetworkHandlerAdapter) healthcheck(resp http.Response, req *http.Request) (interface{}, error) {
	return nil, nil
}
