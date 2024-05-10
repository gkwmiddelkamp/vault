package server

import (
	"vault/server/handlers"
	"vault/vault"
)

const (
	objectIdRegex       = `([a-z0-9]{24})`
	objectIdOrNameRegex = `([a-zA-Z0-9\s-_]+)`
)

func Load(mux *CustomMux) {

	// Public Endpoints
	mux.AddRoute(vault.NewRoute("/", "GET", vault.Public, handlers.IndexHandler))
	// Health checks for application
	mux.AddRoute(vault.NewRoute("/ready", "GET", vault.Public, handlers.ReadinessHandler))
	mux.AddRoute(vault.NewRoute("/live", "GET", vault.Public, handlers.LivenessHandler))

	// Authenticated endpoints
	mux.AddRoute(vault.NewRoute("/token", "GET", vault.EnvironmentAdmin, handlers.TokenListHandler))
	mux.AddRoute(vault.NewRoute("/token", "POST", vault.EnvironmentAdmin, handlers.TokenPostHandler))
	mux.AddRoute(vault.NewRoute("/token/"+objectIdRegex, "GET", vault.EnvironmentAdmin, handlers.TokenGetHandler))
	mux.AddRoute(vault.NewRoute("/token/"+objectIdRegex, "DELETE", vault.EnvironmentAdmin, handlers.TokenDeleteHandler))

	mux.AddRoute(vault.NewRoute("/environment", "GET", vault.MasterAdmin, handlers.EnvironmentListHandler))
	mux.AddRoute(vault.NewRoute("/environment", "POST", vault.MasterAdmin, handlers.EnvironmentPostHandler))
	mux.AddRoute(vault.NewRoute("/environment/"+objectIdRegex, "GET", vault.MasterAdmin, handlers.EnvironmentGetHandler))

	mux.AddRoute(vault.NewRoute("/secret", "GET", vault.ReadWrite, handlers.SecretListHandler))
	mux.AddRoute(vault.NewRoute("/secret", "POST", vault.ReadWrite, handlers.SecretPostHandler))

	mux.AddRoute(vault.NewRoute("/secret/"+objectIdOrNameRegex, "GET", vault.ReadWrite, handlers.SecretGetHandler))
	mux.AddRoute(vault.NewRoute("/secret/"+objectIdOrNameRegex, "DELETE", vault.ReadWrite, handlers.SecretDeleteHandler))
	mux.AddRoute(vault.NewRoute("/secret/"+objectIdOrNameRegex+"/decrypt", "GET", vault.ReadOnly, handlers.SecretDecryptHandler))

}
