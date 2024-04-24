package vault

import (
	"net/http"
	"vault/database"
)

type Interceptor interface {
	Before(w http.ResponseWriter, r *http.Request, cfg *InterceptorConfig) Result
}

type InterceptorConfig struct {
	db            *database.MongoDB
	route         *Route
	handlerConfig *HandlerConfig
}

func (i *InterceptorConfig) GetDb() *database.MongoDB {
	return i.db
}

func (i *InterceptorConfig) GetMinimalAuth() TokenType {
	return i.route.minimalAuth
}

func (i *InterceptorConfig) GetRoute() *Route {
	return i.route
}

func (i *InterceptorConfig) SetHandlerConfig(handlerConfig *HandlerConfig) {
	i.handlerConfig = handlerConfig
}

func (i *InterceptorConfig) GetHandlerConfig() *HandlerConfig {
	return i.handlerConfig
}

func NewInterceptorConfig(db *database.MongoDB, route *Route, handlerConfig *HandlerConfig) *InterceptorConfig {
	return &InterceptorConfig{
		db:            db,
		route:         route,
		handlerConfig: handlerConfig,
	}
}
