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
	handler       Route
	handlerConfig *HandlerConfig
}

func (i *InterceptorConfig) GetDb() *database.MongoDB {
	return i.db
}

func (i *InterceptorConfig) GetMinimalAuth() TokenType {
	return i.handler.minimalAuth
}

func (i *InterceptorConfig) SetHandlerConfig(handlerConfig *HandlerConfig) {
	i.handlerConfig = handlerConfig
}

func (i *InterceptorConfig) GetHandlerConfig() *HandlerConfig {
	return i.handlerConfig
}

func NewInterceptorConfig(db *database.MongoDB, handler Route, handlerConfig *HandlerConfig) *InterceptorConfig {
	return &InterceptorConfig{
		db:            db,
		handler:       handler,
		handlerConfig: handlerConfig,
	}
}
