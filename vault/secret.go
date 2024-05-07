package vault

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const SecretCollection = "secret"

type Secret struct {
	Id             primitive.ObjectID `bson:"_id"`
	Description    string             `bson:"description,omitempty"`
	Secret         string             `bson:"secret,omitempty"`
	EnvironmentId  primitive.ObjectID `bson:"environmentId,omitempty"`
	CreatedAt      primitive.DateTime `bson:"createdAt,omitempty"`
	CreatedBy      string             `bson:"createdBy,omitempty"`
	LastModifiedAt primitive.DateTime `bson:"lastModifiedAt,omitempty"`
	LastModifiedBy string             `bson:"lastModifiedBy,omitempty"`
}

func NewSecret(description string, secret string, environmentId primitive.ObjectID, createdBy string) Secret {
	return Secret{
		Id:            primitive.NewObjectID(),
		Description:   description,
		Secret:        secret,
		EnvironmentId: environmentId,
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		CreatedBy:     createdBy,
	}
}
