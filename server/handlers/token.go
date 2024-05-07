package handlers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"vault/server/models"
	"vault/server/responses"
	"vault/server/services"
	"vault/server/views"
	"vault/vault"
)

type tokenList struct {
}

var TokenListHandler vault.Handler = &tokenList{}

func (h *tokenList) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := bson.D{{"environmentId", cfg.GetEnvironmentId()}}
		if cfg.GetTokenType() == vault.MasterAdmin {
			filter = bson.D{}
		}
		tokenService := services.NewTokenService(cfg.GetDb())
		results, err := tokenService.ListToken(&filter)
		if err != nil {
			response, _ := json.Marshal(responses.NewError(http.StatusNotFound, "Tokens not found"))
			w.Write(response)
			return
		}
		responseList := make([]views.TokenView, 0)
		for _, res := range results {
			responseList = append(responseList, views.NewTokenView(res))
		}
		response, err := json.Marshal(responseList)
		w.Write(response)
	}
}

type tokenGet struct {
}

var TokenGetHandler vault.Handler = &tokenGet{}

func (h *tokenGet) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenService := services.NewTokenService(cfg.GetDb())
		filter, err := tokenService.GetTokenFilter(cfg.GetRouteParam(0), cfg.GetEnvironmentId(), cfg.GetTokenType())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := tokenService.GetToken(filter)
		if err != nil {
			ThrowError(w, http.StatusNotFound, err.Error())
			return
		}
		response, err := json.Marshal(views.NewTokenView(result))
		w.Write(response)
	}
}

type tokenPost struct {
}

var TokenPostHandler vault.Handler = &tokenPost{}

func (h *tokenPost) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.TokenCreate

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			ThrowError(w, http.StatusBadRequest, "Invalid request")
			return
		}

		tokenType, err := vault.TokenTypeFromString(requestBody.TokenType)
		if err != nil {
			ThrowError(w, http.StatusBadRequest, "invalid token type")
			return
		}

		if tokenType > cfg.GetTokenType() {
			ThrowError(w, http.StatusBadRequest, "not allowed to create a token above your level")
			return
		}

		tokenService := services.NewTokenService(cfg.GetDb())
		dbToken, err := tokenService.CreateToken(requestBody, tokenType, cfg.GetTokenDescription(), cfg.GetEnvironmentId(), cfg.GetEnvironmentSecret())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}

		response, _ := json.Marshal(views.NewTokenCreateView(dbToken))
		w.Write(response)
	}
}

type tokenDelete struct {
}

var TokenDeleteHandler vault.Handler = &tokenDelete{}

func (h *tokenDelete) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenService := services.NewTokenService(cfg.GetDb())
		filter, err := tokenService.GetTokenFilter(cfg.GetRouteParam(0), cfg.GetEnvironmentId(), cfg.GetTokenType())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := tokenService.GetToken(filter)
		if err != nil {
			ThrowError(w, http.StatusNotFound, err.Error())
			return
		}

		err = tokenService.DeleteToken(result.Id, cfg.GetEnvironmentId())
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(204)
	}
}
