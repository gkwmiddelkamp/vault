package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"vault/server/models"
	"vault/vault"
)

type tokenList struct {
}

var TokenListHandler vault.Handler = &tokenList{}

func (h *tokenList) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		version := models.Version{Version: "Tokens: 2024.05.0", Name: cfg.GetEnvironmentId().Hex()}
		response, _ := json.Marshal(version)
		log.Println("Writing version: " + version.Version)

		w.Write(response)
	}
}
