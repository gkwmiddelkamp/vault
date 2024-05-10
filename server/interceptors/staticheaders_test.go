package interceptors

import (
	"net/http"
	"testing"
	"vault/vault"
)

func Test_staticHeadersInterceptor_Before(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	testRoute := vault.NewRoute("/", "GET", vault.Public, nil)
	handlerConfig := vault.HandlerConfig{}
	interceptorConfig := vault.NewInterceptorConfig(nil, &testRoute, &handlerConfig)
	responseWriter := MockResponseWriter{}
	staticHeadersInterceptor{}.Before(&responseWriter, request, interceptorConfig)

	headersThatShouldExist := []string{"Access-Control-Allow-Methods", "Content-Type", "X-Content-Type-Options", "X-Frame-Options"}

	for _, v := range headersThatShouldExist {
		if responseWriter.Headers.Get(v) == "" {
			t.Fatalf("staticHeadersInterceptor should set header %s", v)
		}
	}

}
