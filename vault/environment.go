package vault

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const EnvironmentCollection = "environment"

type Environment struct {
	Id             primitive.ObjectID `bson:"_id"`
	Name           string             `bson:"name,omitempty"`
	Contact        string             `bson:"contact,omitempty"`
	Active         bool               `bson:"active,omitempty"`
	CreatedAt      primitive.DateTime `bson:"createdAt,omitempty"`
	CreatedBy      string             `bson:"createdBy,omitempty"`
	LastModifiedAt primitive.DateTime `bson:"lastModifiedAt,omitempty"`
	LastModifiedBy string             `bson:"lastModifiedBy,omitempty"`
}

func NewEnvironment(name string, contact string, active bool, createdBy string) Environment {
	return Environment{
		Id:        primitive.NewObjectID(),
		Name:      name,
		Contact:   contact,
		Active:    active,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		CreatedBy: createdBy,
	}
}
