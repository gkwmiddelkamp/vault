package handlers

import (
	"net/http"
	"vault/vault"
)

type empty struct {
}

var EmptyHandler vault.Handler = &empty{}

func (h *empty) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
