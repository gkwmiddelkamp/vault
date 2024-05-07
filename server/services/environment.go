package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"vault/database"
	"vault/internal"
	"vault/server/models"
	"vault/vault"
)

type EnvironmentService struct {
	db *database.MongoDB
}

func (s *EnvironmentService) ListEnvironments() ([]*vault.Environment, error) {
	var results []*vault.Environment
	filter := bson.D{{}}

	cursor, err := s.db.Connection.Collection(vault.EnvironmentCollection).Find(context.Background(), filter)
	if err != nil {
		log.Println("No environments found: " + err.Error())
		return nil, errors.New("no environments found")
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println("Results could not be parsed: " + err.Error())
		return nil, errors.New("results could not be parsed")
	}
	return results, nil
}

func (s *EnvironmentService) GetEnvironment(id string) (*vault.Environment, error) {
	var result vault.Environment
	requestedId, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", requestedId}}

	if err = s.db.Connection.Collection(vault.EnvironmentCollection).FindOne(context.Background(), filter).Decode(&result); err != nil {
		log.Println("Environment not found: " + err.Error())
		return nil, errors.New("environment not found")
	}
	return &result, nil
}

func (s *EnvironmentService) CreateEnvironment(create models.EnvironmentCreate, createdBy string) (*vault.Environment, *[]byte, error) {
	environmentSecret := internal.RandomString(vault.KeyLength)
	environment := vault.NewEnvironment(create.Name, create.Contact, create.Active, createdBy)
	insertResult, err := s.db.Connection.Collection(vault.EnvironmentCollection).InsertOne(context.Background(), environment)
	if err != nil {
		return nil, nil, err
	}

	environmentId := insertResult.InsertedID.(primitive.ObjectID)

	insertedEnvironment, err := s.GetEnvironment(environmentId.Hex())
	if err != nil {
		return nil, nil, err
	}
	return insertedEnvironment, &environmentSecret, nil

}

func (s *EnvironmentService) DeleteEnvironment(id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}
	deleteResult, err := s.db.Connection.Collection(vault.EnvironmentCollection).DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 1 {
		return nil
	}
	return errors.New("invalid request")

}

func NewEnvironmentService(db *database.MongoDB) EnvironmentService {
	return EnvironmentService{
		db: db,
	}
}
