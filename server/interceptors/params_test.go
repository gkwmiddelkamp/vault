package interceptors

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"vault/vault"
)

type MockResponseWriter struct {
	Body    io.ReadCloser
	Code    int
	Headers http.Header
}

func (e *MockResponseWriter) Header() http.Header {
	if e.Headers == nil {
		e.Headers = make(http.Header)
	}
	return e.Headers
}

func (e *MockResponseWriter) Write(data []byte) (int, error) {
	e.Body = io.NopCloser(bytes.NewReader(data))
	return 0, fmt.Errorf("always errors")
}

func (e *MockResponseWriter) WriteHeader(statusCode int) {
	e.Code = statusCode
}

func Test_paramsInterceptor_Before(t *testing.T) {
	request, _ := http.NewRequest("GET", "/test/test123A/anothertest", nil)
	testRoute := vault.NewRoute("/test/([a-z0-9A-Z]+)/([a-z]+)", "GET", vault.Public, nil)
	handlerConfig := vault.HandlerConfig{}
	interceptorConfig := vault.NewInterceptorConfig(nil, &testRoute, &handlerConfig)
	responseWriter := MockResponseWriter{}
	paramsInterceptor{}.Before(&responseWriter, request, interceptorConfig)

	if handlerConfig.GetRouteParam(0) != "test123A" || handlerConfig.GetRouteParam(1) != "anothertest" {
		t.Fatal("ParamsInterceptor should put route params when they are requested")
	}

}
