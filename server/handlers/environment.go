package handlers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"vault/server/models"
	"vault/server/responses"
	"vault/server/services"
	"vault/server/views"
	"vault/vault"
)

type environmentList struct {
}

var EnvironmentListHandler vault.Handler = &environmentList{}

// EnvironmentList godoc
//
//	@Summary		Environment List
//	@Description	Get a list of environments
//	@Tags			Environment
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	views.EnvironmentView
//	@Failure		404	{string}	string	"Not found"
//	@Router			/environment [get]
func (h *environmentList) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := services.NewEnvironmentService(cfg.GetDb())
		results, err := service.ListEnvironments()
		if err != nil {
			response, _ := json.Marshal(responses.NewError(http.StatusNotFound, "Error while fetching environments"))
			w.Write(response)
			return
		}
		responseList := make([]views.EnvironmentView, 0)
		for _, res := range results {
			responseList = append(responseList, views.NewEnvironmentView(res))
		}
		response, err := json.Marshal(responseList)
		w.Write(response)
	}
}

type environmentGet struct {
}

var EnvironmentGetHandler vault.Handler = &environmentGet{}

// EnvironmentGet godoc
//
//	@Summary		Environment Get
//	@Description	Get a single environment
//	@Tags			Environment
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	views.EnvironmentView
//	@Failure		404	{string}	string	"Not found"
//	@Router			/environment/{id} [get]
func (h *environmentGet) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := services.NewEnvironmentService(cfg.GetDb())
		result, err := service.GetEnvironment(cfg.GetRouteParam(0))
		if err != nil {
			response, _ := json.Marshal(responses.NewError(http.StatusNotFound, "Environment not found"))
			w.Write(response)
			return
		}

		response, err := json.Marshal(views.NewEnvironmentView(result))
		w.Write(response)
	}
}

type environmentPost struct {
}

var EnvironmentPostHandler vault.Handler = &environmentPost{}

// EnvironmentPost godoc
//
// @Summary		Environment Create
// @Description	Create a new environment
// @Tags			Environment
// @Accept			json
// @Produce		json
// @Success		200	{object}	views.EnvironmentCreateView
// @Failure		401	{string}	string	"Unauthorized"
// @Failure		402	{string}	string	"Bad request"
// @Router			/environment [post]
// @Param 			request body models.EnvironmentCreate true "Environment create object"
func (h *environmentPost) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.EnvironmentCreate
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			ThrowError(w, http.StatusBadRequest, "Invalid request")
			return
		}

		environmentService := services.NewEnvironmentService(cfg.GetDb())
		dbEnvironment, environmentSecret, err := environmentService.CreateEnvironment(requestBody, cfg.GetTokenDescription())
		if err != nil {
			log.Println("Environment not created: " + err.Error())
			ThrowError(w, http.StatusBadRequest, "environment not created")
			return
		}

		tokenService := services.NewTokenService(cfg.GetDb())
		tokenCreate := models.TokenCreate{Description: requestBody.Name}
		token, err := tokenService.CreateToken(tokenCreate, vault.EnvironmentAdmin, cfg.GetTokenDescription(), dbEnvironment.Id, *environmentSecret)
		if err != nil {
			log.Println("Environment partially created: " + err.Error())
			ThrowError(w, http.StatusBadRequest, "environment partially created")
			return
		}
		response, err := json.Marshal(views.NewEnvironmentCreateView(dbEnvironment, token.GetSecret()))
		w.Write(response)
	}
}

type environmentDelete struct {
}

var EnvironmentDeleteHandler vault.Handler = &environmentDelete{}

// EnvironmentDelete godoc
//
// @Summary		Environment Delete
// @Description	Delete an existing environment and all objects belonging to that environment
// @Tags			Environment
// @Accept			json
// @Produce		json
// @Success		201	{string}	string "Accepted"
// @Failure		401	{string}	string	"Unauthorized"
// @Failure		402	{string}	string	"Bad request"
// @Router			/environment/{id} [delete]
func (h *environmentDelete) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		environmentService := services.NewEnvironmentService(cfg.GetDb())
		environment, err := environmentService.GetEnvironment(cfg.GetRouteParam(0))
		if err != nil {
			ThrowError(w, http.StatusNotFound, err.Error())
			return
		}

		tokenService := services.NewTokenService(cfg.GetDb())
		filter := bson.D{{"environmentId", environment.Id}}
		tokens, err := tokenService.ListToken(&filter)
		for _, token := range tokens {
			err = tokenService.DeleteToken(token.Id, environment.Id)
			if err != nil {
				ThrowError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		err = environmentService.DeleteEnvironment(environment.Id)
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		//TODO: handle delete passwords

	}
}
