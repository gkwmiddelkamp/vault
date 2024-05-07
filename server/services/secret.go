package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"regexp"
	"vault/database"
	"vault/server/models"
	"vault/vault"
)

const descriptionRegex = `^([a-zA-Z0-9\s-_]+)$`

type SecretService struct {
	db *database.MongoDB
}

func NewSecretService(db *database.MongoDB) SecretService {
	return SecretService{db: db}
}

func (s *SecretService) GetSecretFilter(id string, environmentId primitive.ObjectID) *bson.D {
	requestedId, err := primitive.ObjectIDFromHex(id)
	var filter bson.D
	if err == nil {
		filter = bson.D{{"_id", requestedId}, {"environmentId", environmentId}}
	} else {
		filter = bson.D{{"description", id}, {"environmentId", environmentId}}
	}

	return &filter
}

func (s *SecretService) ListSecret(environmentId primitive.ObjectID) ([]*vault.Secret, error) {
	var results []*vault.Secret
	filter := bson.D{{"environmentId", environmentId}}
	opts := options.Find().SetProjection(bson.D{{"secret", 0}})
	cursor, err := s.db.Connection.Collection(vault.SecretCollection).Find(context.Background(), filter, opts)
	if err != nil {
		return nil, errors.New("no secrets found")
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, errors.New("could not map secrets")
	}
	return results, nil
}

func (s *SecretService) GetSecret(filter *bson.D, includeSecret bool) (*vault.Secret, error) {
	var result vault.Secret
	opts := options.FindOne().SetProjection(bson.D{{"secret", 0}})
	if includeSecret {
		opts = options.FindOne()
	}
	if err := s.db.Connection.Collection(vault.SecretCollection).FindOne(context.Background(), filter, opts).Decode(&result); err != nil {
		return nil, errors.New("secret not found")
	}
	return &result, nil
}

func (s *SecretService) CreateSecret(create models.SecretCreate, createdBy string, environmentId primitive.ObjectID, environmentSecret []byte) (*vault.Secret, error) {
	regex, _ := regexp.Compile(descriptionRegex)
	if !regex.MatchString(create.Description) {
		return nil, errors.New("invalid description, should match: " + descriptionRegex)
	}
	alreadyExistsFilter := s.GetSecretFilter(create.Description, environmentId)
	if _, err := s.GetSecret(alreadyExistsFilter, false); err == nil {
		return nil, errors.New("a secret with this name already exists")
	}

	dbSecret := vault.NewSecret(create.Description, create.Secret, environmentId, createdBy)
	// encrypt password

	_, err := s.db.Connection.Collection(vault.SecretCollection).InsertOne(context.Background(), dbSecret)
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil, errors.New("failed to save the new token")
	}

	return &dbSecret, err
}

func (s *SecretService) DeleteSecret(id primitive.ObjectID, environmentId primitive.ObjectID) error {
	// Can only delete tokens in their own environment
	filter := bson.D{{"_id", id}, {"environmentId", environmentId}}
	deleteResult, err := s.db.Connection.Collection(vault.SecretCollection).DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 1 {
		return nil
	}
	return errors.New("invalid request")
}
