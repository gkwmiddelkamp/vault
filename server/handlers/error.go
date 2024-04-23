package handlers

import (
	"encoding/json"
	"net/http"
	"vault/server/responses"
	"vault/vault"
)

type error struct {
}

var ErrorHandler vault.Handler = &error{}

func (h *error) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(responses.NewError(403, "No way"))

		_, err := w.Write(response)
		if err != nil {
			return
		}
	}
}
