package vault

import (
	"net/http"
)

type Route struct {
	minimalAuth TokenType
	method      string
	pattern     string
	handler     Handler
}

func NewRoute(pattern string, method string, minimalAuth TokenType, handler Handler) Route {
	return Route{
		pattern:     pattern,
		method:      method,
		minimalAuth: minimalAuth,
		handler:     handler,
	}
}

func (r *Route) GetPattern() string {
	return r.pattern
}

func (r *Route) GetMethod() string {
	return r.method
}

func (r *Route) Handle(cfg HandlerConfig) http.HandlerFunc {
	return r.handler.Handle(cfg)
}
