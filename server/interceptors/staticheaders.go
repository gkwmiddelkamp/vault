package interceptors

import (
	"net/http"
	"vault/vault"
)

type staticHeadersInterceptor struct{}

var StaticHeadersInterceptor vault.Interceptor = staticHeadersInterceptor{}

func (i staticHeadersInterceptor) Before(w http.ResponseWriter, r *http.Request, cfg *vault.InterceptorConfig) vault.Result {
	h := w.Header()

	h.Set("Content-Type", "application/json")
	// Set HSTS
	h.Set("Access-Control-Allow-Methods", "GET,POST,PUT")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "Deny")

	return vault.NotDone()
}
