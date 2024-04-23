package handlers

import (
	"github.com/heptiolabs/healthcheck"
	"net/http"
	"time"
	"vault/vault"
)

type liveness struct {
}

var LivenessHandler vault.Handler = &liveness{}

func (h *liveness) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := healthcheck.NewHandler()
		health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
		health.ReadyEndpoint(w, r)
	}
}

type readiness struct {
}

var ReadinessHandler vault.Handler = &readiness{}

func (h *readiness) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := healthcheck.NewHandler()
		health.AddReadinessCheck(
			"upstream-dep-dns",
			healthcheck.DNSResolveCheck("github.com", 50*time.Millisecond))
		//health.AddReadinessCheck("mongodb-connection", healthcheck.Async(checkMongoDbConnection(cfg.GetDb()), 3000))
		health.ReadyEndpoint(w, r)

	}
}
