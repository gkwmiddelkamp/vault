package server

import (
	"log"
	"net/http"
	"regexp"
	"vault/database"
	"vault/server/interceptors"
	"vault/vault"
)

type CustomMux struct {
	defaultMux   *http.ServeMux
	routes       []vault.Route
	interceptors []vault.Interceptor
	db           *database.MongoDB
}

func NewCustomMux(db *database.MongoDB) CustomMux {
	return CustomMux{db: db}
}

// Ensure http.Handler interface is satisfied
func (m *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.defaultMux == nil {
		m.init()
	}

	log.Printf("Requesting %s", r.URL.Path)

	cfg := vault.NewHandlerConfig(m.db)
	for _, route := range m.routes {

		if r.Method == route.GetMethod() {
			regex, _ := regexp.Compile(route.GetPattern() + "$")

			if regex.MatchString(r.URL.Path) {
				for _, i := range m.interceptors {
					res := i.Before(w, r, vault.NewInterceptorConfig(m.db, &route, &cfg))
					if res.Done {
						return
					}
				}
				log.Println(route.GetMethod() + " " + route.GetPattern())

				route.Handle(cfg).ServeHTTP(w, r)
				return
			}
		}
	}

	http.NotFound(w, r)
}

func (m *CustomMux) AddRoute(route vault.Route) {
	m.routes = append(m.routes, route)
}

func (m *CustomMux) init() {
	m.defaultMux = http.NewServeMux()

	m.interceptors = append(m.interceptors, interceptors.StaticHeadersInterceptor)
	m.interceptors = append(m.interceptors, interceptors.AuthInterceptor)
	m.interceptors = append(m.interceptors, interceptors.ParamsInterceptor)

}
