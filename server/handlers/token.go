package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
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
