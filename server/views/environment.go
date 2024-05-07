package views

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vault/vault"
)

type EnvironmentView struct {
	Id             primitive.ObjectID `json:"id"`
	Name           string             `json:"name,omitempty"`
	Contact        string             `json:"contact,omitempty"`
	Active         bool               `json:"active,omitempty"`
	CreatedAt      string             `json:"createdAt,omitempty"`
	CreatedBy      string             `json:"createdBy,omitempty"`
	LastModifiedAt string             `json:"lastModifiedAt,omitempty"`
	LastModifiedBy string             `json:"lastModifiedBy,omitempty"`
}

func NewEnvironmentView(environment *vault.Environment) EnvironmentView {
	result := EnvironmentView{
		Id:        environment.Id,
		Name:      environment.Name,
		Contact:   environment.Contact,
		Active:    environment.Active,
		CreatedAt: environment.CreatedAt.Time().String(),
		CreatedBy: environment.CreatedBy,
	}

	if environment.LastModifiedBy != "" {
		result.LastModifiedAt = environment.LastModifiedAt.Time().String()
		result.LastModifiedBy = environment.LastModifiedBy
	}

	return result
}
