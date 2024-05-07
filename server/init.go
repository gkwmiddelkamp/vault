package server

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"vault/database"
	"vault/internal"
	"vault/vault"
)

func InitDatabase(db *database.MongoDB) {
	count, err := db.Connection.Collection(vault.EnvironmentCollection).CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("An error occurred during DB initialization: " + err.Error())
	}
	if count > 0 {
		log.Println("Database already initialized")
		return
	}
	log.Println("Empty dataset found, initializing")
	environmentSecret := internal.RandomString(vault.KeyLength)
	environment := vault.NewEnvironment("Master environment", "no-reply@vault.pnck.nl", true, "System")
	one, err := db.Connection.Collection(vault.EnvironmentCollection).InsertOne(context.Background(), environment)
	if err != nil {
		log.Fatal("No documents inserted")
	}

	environmentId := one.InsertedID.(primitive.ObjectID)

	token := vault.NewToken("Master token", environmentId, "System", vault.MasterAdmin, environmentSecret)
	db.Connection.Collection(vault.TokenCollection).InsertOne(context.Background(), token)
	log.Println("Initialization done, your master token is: " + token.GetSecret())
}
