package views

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vault/vault"
)

type SecretView struct {
	Id             primitive.ObjectID `json:"id"`
	Description    string             `json:"description,omitempty"`
	CreatedAt      string             `json:"createdAt,omitempty"`
	CreatedBy      string             `json:"createdBy,omitempty"`
	LastModifiedAt string             `json:"lastModifiedAt,omitempty"`
	LastModifiedBy string             `json:"lastModifiedBy,omitempty"`
}

func NewSecretView(secret *vault.Secret) SecretView {
	result := SecretView{
		Id:          secret.Id,
		Description: secret.Description,
		CreatedAt:   secret.CreatedAt.Time().String(),
		CreatedBy:   secret.CreatedBy,
	}
	if secret.LastModifiedAt > 0 {
		result.LastModifiedAt = secret.LastModifiedAt.Time().String()
		result.LastModifiedBy = secret.LastModifiedBy
	}
	return result
}
