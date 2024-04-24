package interceptors

import (
	"net/http"
	"regexp"
	"vault/vault"
)

type paramsInterceptor struct{}

var ParamsInterceptor vault.Interceptor = paramsInterceptor{}

func (i paramsInterceptor) Before(w http.ResponseWriter, r *http.Request, cfg *vault.InterceptorConfig) vault.Result {
	regex, _ := regexp.Compile(cfg.GetRoute().GetPattern() + "$")
	results := regex.FindStringSubmatch(r.URL.Path)
	var params []string
	if len(results) > 1 {
		for _, result := range results[1:] {
			params = append(params, result)
		}
	}
	cfg.GetHandlerConfig().SetRouteParams(params)
	return vault.NotDone()
}
