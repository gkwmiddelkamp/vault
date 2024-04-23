package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"vault/server/models"
	"vault/vault"
)

type index struct {
}

var IndexHandler vault.Handler = &index{}

func (h *index) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version := models.Version{Version: "Version: 2024.05.0", Name: cfg.GetDb().Connection.Name()}
		response, _ := json.Marshal(version)
		log.Println("Writing version: " + version.Version)

		w.Write(response)
	}
}
