package handlers

import (
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"os"
	"vault/vault"
)

type swaggerRedirect struct {
}

var SwaggerRedirectHandler vault.Handler = &swaggerRedirect{}

func (h *swaggerRedirect) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/docs/index.html")
		w.WriteHeader(302)
	}
}

type swaggerBase struct {
}

var SwaggerBaseHandler vault.Handler = &swaggerBase{}

func (h *swaggerBase) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	)
}

type swaggerJson struct {
}

var SwaggerJsonHandler vault.Handler = &swaggerJson{}

func (h *swaggerJson) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swaggerJson, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			ThrowError(w, 404, err.Error())
		}
		w.Write(swaggerJson)
	}
}
