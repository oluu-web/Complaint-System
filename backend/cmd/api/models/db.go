package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

var client *mongo.Client
var dbName = os.Getenv("DB")

func ConnectToDB() error {

	mongoURI := os.Getenv("MONGOURI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %w", err)
	}
	log.Println("connected successfully")
	return nil
}

// GetDBCollection returns a reference to a collection in a MongoDB database
func GetDBCollection(collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
