package views

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vault/vault"
)

type TokenView struct {
	Id            primitive.ObjectID `json:"id"`
	Description   string             `json:"description,omitempty"`
	EnvironmentId primitive.ObjectID `json:"environmentId,omitempty"`
	CreatedAt     string             `json:"createdAt,omitempty"`
	CreatedBy     string             `json:"createdBy,omitempty"`
	ExpiresAt     string             `json:"expiresAt,omitempty"`
	TokenType     string             `json:"tokenType,omitempty"`
}

func NewTokenView(token *vault.Token) TokenView {
	result := TokenView{
		Id:            token.Id,
		Description:   token.Description,
		EnvironmentId: token.EnvironmentId,
		CreatedAt:     token.CreatedAt.Time().String(),
		CreatedBy:     token.CreatedBy,
		TokenType:     token.TokenType.String(),
	}
	if token.ExpiresAt > 0 {
		result.ExpiresAt = token.ExpiresAt.Time().String()
	}
	return result
}
