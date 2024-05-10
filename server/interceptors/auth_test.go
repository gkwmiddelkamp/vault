package interceptors

import (
	"net/http"
	"testing"
	"vault/vault"
)

func Test_authInterceptor_Before(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	testRoute := vault.NewRoute("/", "GET", vault.EnvironmentAdmin, nil)
	handlerConfig := vault.HandlerConfig{}
	interceptorConfig := vault.NewInterceptorConfig(nil, &testRoute, &handlerConfig)
	responseWriter := MockResponseWriter{}
	authInterceptor{}.Before(&responseWriter, request, interceptorConfig)

	if responseWriter.Code != 401 {
		t.Fatal("AuthInterceptor should throw unauthorized when no token is sent")
	}
	responseWriter.WriteHeader(0)

	request.Header.Add("X-API-Key", "test")
	authInterceptor{}.Before(&responseWriter, request, interceptorConfig)
	if responseWriter.Code != 401 {
		t.Fatal("AuthInterceptor should throw unauthorized when invalid token is sent")
	}
}
