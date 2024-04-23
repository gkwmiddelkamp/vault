package vault

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"vault/database"
)

type HandlerConfig struct {
	db            *database.MongoDB
	tokenType     TokenType
	environmentId primitive.ObjectID
}

func NewHandlerConfig(db *database.MongoDB) HandlerConfig {
	return HandlerConfig{
		db: db,
	}
}

type Handler interface {
	Handle(cfg HandlerConfig) http.HandlerFunc
}

func (h *HandlerConfig) GetDb() *database.MongoDB {
	return h.db
}

func (h *HandlerConfig) GetTokenType() TokenType {
	return h.tokenType
}

func (h *HandlerConfig) SetTokenType(t TokenType) {
	h.tokenType = t
}

func (h *HandlerConfig) GetEnvironmentId() primitive.ObjectID {
	return h.environmentId
}

func (h *HandlerConfig) SetEnvironmentId(id primitive.ObjectID) {
	h.environmentId = id
}
