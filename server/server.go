package server

import (
	"vault/server/handlers"
	"vault/vault"
)

const (
	objectIdRegex = "([a-z0-9]{24})"
)

func Load(mux *CustomMux) {

	// Public Endpoints
	mux.AddRoute(vault.NewRoute("/favicon.ico", "GET", vault.Public, handlers.EmptyHandler))
	mux.AddRoute(vault.NewRoute("/error", "GET", vault.Public, handlers.ErrorHandler))
	mux.AddRoute(vault.NewRoute("/", "GET", vault.Public, handlers.IndexHandler))
	// Health checks for application
	mux.AddRoute(vault.NewRoute("/ready", "GET", vault.Public, handlers.ReadinessHandler))
	mux.AddRoute(vault.NewRoute("/live", "GET", vault.Public, handlers.LivenessHandler))

	// Authenticated endpoints
	mux.AddRoute(vault.NewRoute("/token", "GET", vault.EnvironmentAdmin, handlers.TokenListHandler))
	mux.AddRoute(vault.NewRoute("/token", "POST", vault.EnvironmentAdmin, handlers.TokenPostHandler))

	mux.AddRoute(vault.NewRoute("/token/"+objectIdRegex, "GET", vault.EnvironmentAdmin, handlers.TokenGetHandler))

}
