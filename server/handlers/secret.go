package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"vault/server/models"
	"vault/server/responses"
	"vault/server/services"
	"vault/server/views"
	"vault/vault"
)

type secretList struct {
}

var SecretListHandler vault.Handler = &secretList{}

func checkTokenType(tokenType vault.TokenType) error {
	if tokenType >= vault.EnvironmentAdmin {
		return errors.New("token leven too high for secret management")
	}
	return nil
}

// SecretListHandler godoc
//
//	@Summary		Secret List
//	@Description	Get a secret list
//	@Tags			Secret
//	@Produce		json
//	@Success		200	{array}	views.SecretView
//	@Router			/secret [get]
func (h *secretList) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkTokenType(cfg.GetTokenType()); err != nil {
			ThrowError(w, http.StatusUnauthorized, err.Error())
			return
		}
		secretService := services.NewSecretService(cfg.GetDb())
		results, err := secretService.ListSecret(cfg.GetEnvironmentId())
		if err != nil {
			response, _ := json.Marshal(responses.NewError(http.StatusNotFound, "Secrets not found"))
			w.Write(response)
			return
		}
		responseList := make([]views.SecretView, 0)

		for _, res := range results {
			responseList = append(responseList, views.NewSecretView(res))
		}

		response, err := json.Marshal(responseList)
		w.Write(response)
	}
}

type secretGet struct {
}

var SecretGetHandler vault.Handler = &secretGet{}

// SecretGetHandler godoc
//
//	@Summary		Secret Get
//	@Description	Get a single secret
//	@Tags			Secret
//	@Produce		json
//	@Success		200	{object}	views.SecretView
//	@Router			/secret/{id} [get]
func (h *secretGet) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkTokenType(cfg.GetTokenType()); err != nil {
			ThrowError(w, http.StatusUnauthorized, err.Error())
			return
		}
		secretService := services.NewSecretService(cfg.GetDb())
		filter := secretService.GetSecretFilter(cfg.GetRouteParam(0), cfg.GetEnvironmentId())

		result, err := secretService.GetSecret(filter, false)
		if err != nil {
			ThrowError(w, http.StatusNotFound, err.Error())
			return
		}
		response, err := json.Marshal(views.NewSecretView(result))
		w.Write(response)
	}
}

type secretPost struct {
}

var SecretPostHandler vault.Handler = &secretPost{}

// SecretPostHandler godoc
//
// @Summary		Secret Create
// @Description	Create a new secret
// @Tags			Secret
// @Accept		json
// @Produce		json
// @Success		200	{object}	views.SecretView
// @Router			/secret [post]
// @Param 	request body models.SecretCreate true "Secret create object"
func (h *secretPost) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkTokenType(cfg.GetTokenType()); err != nil {
			ThrowError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var requestBody models.SecretCreate

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			ThrowError(w, http.StatusBadRequest, "Invalid request")
			return
		}

		secretService := services.NewSecretService(cfg.GetDb())
		dbSecret, err := secretService.CreateSecret(requestBody, cfg.GetTokenDescription(), cfg.GetEnvironmentId(), cfg.GetEnvironmentSecret())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}

		response, _ := json.Marshal(views.NewSecretView(dbSecret))
		w.Write(response)
	}
}

type secretDelete struct {
}

var SecretDeleteHandler vault.Handler = &secretDelete{}

// SecretDeleteHandler godoc
//
// @Summary		Secret Delete
// @Description	Delete a secret
// @Tags			Secret
// @Produce		json
// @Success		204	{nil}  string "Accepted"
// @Router			/secret/{id} [delete]
func (h *secretDelete) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkTokenType(cfg.GetTokenType()); err != nil {
			ThrowError(w, http.StatusUnauthorized, err.Error())
			return
		}
		secretService := services.NewSecretService(cfg.GetDb())
		filter := secretService.GetSecretFilter(cfg.GetRouteParam(0), cfg.GetEnvironmentId())

		result, err := secretService.GetSecret(filter, false)
		if err != nil {
			ThrowError(w, http.StatusNotFound, err.Error())
			return
		}

		err = secretService.DeleteSecret(result.Id, cfg.GetEnvironmentId())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(204)
	}
}

type secretDecrypt struct {
}

var SecretDecryptHandler vault.Handler = &secretDecrypt{}

// SecretDecryptHandler godoc
//
// @Summary		Secret Decrypt
// @Description	Get a decrypted secret
// @Tags			Secret
// @Produce		json
// @Success		200	{object} views.SecretDecryptView
// @Router			/secret/{id}/decrypt [get]
func (h *secretDecrypt) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkTokenType(cfg.GetTokenType()); err != nil {
			ThrowError(w, http.StatusUnauthorized, err.Error())
			return
		}
		secretService := services.NewSecretService(cfg.GetDb())
		filter := secretService.GetSecretFilter(cfg.GetRouteParam(0), cfg.GetEnvironmentId())

		result, err := secretService.GetSecret(filter, true)
		if err != nil {
			ThrowError(w, http.StatusNotFound, err.Error())
			return
		}

		decryptedSecret, err := secretService.DecryptSecret(result.Secret, cfg.GetEnvironmentSecret())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		response, err := json.Marshal(views.NewSecretDecodeView(string(decryptedSecret)))
		w.Write(response)
	}
}
