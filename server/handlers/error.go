package handlers

import (
	"encoding/json"
	"net/http"
	"vault/server/responses"
)

func ThrowError(w http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(responses.NewError(code, message))
	w.WriteHeader(code)
	w.Write(response)
}
