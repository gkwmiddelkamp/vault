package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
	"vault/server/models"
	"vault/server/responses"
	"vault/server/views"
	"vault/vault"
)

type tokenList struct {
}

var TokenListHandler vault.Handler = &tokenList{}

func (h *tokenList) Handle(cfg vault.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var results []vault.Token
		filter := bson.D{{"environmentId", cfg.GetEnvironmentId()}}
		if cfg.GetTokenType() == vault.MasterAdmin {
			filter = bson.D{}
		}
		log.Println(cfg.GetEnvironmentId())
		cursor, err := cfg.GetDb().Connection.Collection(vault.TokenCollection).Find(context.Background(), filter)
		if err != nil {
			log.Println("test" + err.Error())
		}
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Println("test2 " + err.Error())
		}

		var responseList []views.TokenView
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
		var result vault.Token
		requestedId, err := primitive.ObjectIDFromHex(cfg.GetRouteParam(0))
		filter := bson.D{{"_id", requestedId}, {"environmentId", cfg.GetEnvironmentId()}}
		if cfg.GetTokenType() == vault.MasterAdmin {
			filter = bson.D{{"_id", requestedId}}
		}

		err = cfg.GetDb().Connection.Collection(vault.TokenCollection).FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			response, _ := json.Marshal(responses.NewError(http.StatusNotFound, "Token not found"))
			w.Write(response)
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
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			ThrowError(w, http.StatusBadRequest, "Invalid request")
			return
		}

		environmentId := cfg.GetEnvironmentId()
		if cfg.GetTokenType() == vault.MasterAdmin && requestBody.EnvironmentId != "" {
			environmentId, err = primitive.ObjectIDFromHex(requestBody.EnvironmentId)
			if err != nil {
				ThrowError(w, http.StatusBadRequest, "Invalid environmentId")
				return
			}
		}

		tokenType, err := vault.TokenTypeFromString(requestBody.TokenType)
		if err != nil {
			ThrowError(w, http.StatusBadRequest, err.Error())
			return
		}
		dbToken := vault.NewToken(requestBody.Description, environmentId, cfg.GetTokenDescription(), tokenType)
		if requestBody.ExpiresAt != "" {
			const layout1 = "2006-01-02"
			const layout2 = "2006-01-02 15:04:05"
			const layout3 = "2006-01-02T15:04:05"

			expiresAt, err := time.Parse(layout1, requestBody.ExpiresAt)
			if err != nil {
				expiresAt, err = time.Parse(layout2, requestBody.ExpiresAt)
				if err != nil {
					expiresAt, err = time.Parse(layout3, requestBody.ExpiresAt)
					if err != nil {
						expiresAt, err = time.Parse(layout2, requestBody.ExpiresAt)
						ThrowError(w, http.StatusBadRequest, "invalid date")
						return
					}
				}
			}

			dbToken.SetExpiresAt(primitive.NewDateTimeFromTime(expiresAt))
		}
		_, err = cfg.GetDb().Connection.Collection(vault.TokenCollection).InsertOne(context.Background(), dbToken)
		if err != nil {
			log.Println("Error: " + err.Error())
			ThrowError(w, http.StatusBadRequest, "Failed to save the new token")
			return
		}

		response, _ := json.Marshal(views.NewTokenCreateView(dbToken))
		w.Write(response)
	}
}
