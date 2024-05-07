package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
	"vault/database"
	"vault/server/models"
	"vault/vault"
)

type TokenService struct {
	db *database.MongoDB
}

func (s *TokenService) GetTokenFilter(id string, environmentId primitive.ObjectID, tokenType vault.TokenType) (*bson.D, error) {
	requestedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	filter := bson.D{{"_id", requestedId}, {"environmentId", environmentId}}
	if tokenType == vault.MasterAdmin {
		filter = bson.D{{"_id", requestedId}}
	}
	return &filter, nil
}

func (s *TokenService) ListToken(filter *bson.D) ([]*vault.Token, error) {
	var results []*vault.Token

	cursor, err := s.db.Connection.Collection(vault.TokenCollection).Find(context.Background(), filter)
	if err != nil {
		return nil, errors.New("no tokens found")
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, errors.New("could not map tokens")
	}
	return results, nil
}

func (s *TokenService) GetToken(filter *bson.D) (*vault.Token, error) {
	var result vault.Token
	err := s.db.Connection.Collection(vault.TokenCollection).FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, errors.New("token not found")
	}
	return &result, nil
}

func (s *TokenService) CreateToken(create models.TokenCreate, tokenType vault.TokenType, createdBy string, environmentId primitive.ObjectID, environmentSecret []byte) (*vault.Token, error) {
	dbToken := vault.NewToken(create.Description, environmentId, createdBy, tokenType, environmentSecret)
	if create.ExpiresAt != "" {
		const layout1 = "2006-01-02"
		const layout2 = "2006-01-02 15:04:05"
		const layout3 = "2006-01-02T15:04:05"

		expiresAt, err := time.Parse(layout1, create.ExpiresAt)
		if err != nil {
			expiresAt, err = time.Parse(layout2, create.ExpiresAt)
			if err != nil {
				expiresAt, err = time.Parse(layout3, create.ExpiresAt)
				if err != nil {
					expiresAt, err = time.Parse(layout2, create.ExpiresAt)
					return nil, errors.New("invalid date")
				}
			}
		}

		dbToken.SetExpiresAt(primitive.NewDateTimeFromTime(expiresAt))
	}
	_, err := s.db.Connection.Collection(vault.TokenCollection).InsertOne(context.Background(), dbToken)
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil, errors.New("failed to save the new token")
	}
	return &dbToken, nil
}

func (s *TokenService) DeleteToken(id primitive.ObjectID, environmentId primitive.ObjectID) error {
	// Can only delete tokens in their own environment
	filter := bson.D{{"_id", id}, {"environmentId", environmentId}}
	deleteResult, err := s.db.Connection.Collection(vault.TokenCollection).DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 1 {
		return nil
	}
	return errors.New("invalid request")

}

func NewTokenService(db *database.MongoDB) TokenService {
	return TokenService{
		db: db,
	}
}
