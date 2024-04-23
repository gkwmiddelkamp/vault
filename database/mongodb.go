package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"log"
	"os"
)

type MongoDB struct {
	Connection *mongo.Database
}

func NewMongoDB(connectionUri string) MongoDB {
	m := MongoDB{}
	m.connect(connectionUri)
	return m
}

func (m *MongoDB) connect(connectionUri string) {
	if connectionUri == "" {
		log.Fatal("No MongoDB uri defined, see documentation for usage")
	}

	databaseName := ""
	cs, err := connstring.ParseAndValidate(connectionUri)
	if err != nil {
		if os.Getenv("MONGODB_DATABASE") == "" {
			log.Fatal("No MongoDB database defined, see documentation for usage")
		}
		databaseName = os.Getenv("MONGODB_DATABASE")
	} else {
		databaseName = cs.Database
	}

	if databaseName == "" {
		log.Fatal("An error occured while searching database name")
	}

	opts := options.Client().ApplyURI(connectionUri)
	client, err := mongo.Connect(context.Background(), opts)

	var result bson.M
	if err := client.Database(databaseName).RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to database successful")

	m.Connection = client.Database(databaseName)
}
