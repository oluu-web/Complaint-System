package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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

func ConnectToDB() error {
	dbName := os.Getenv("DB_NAME")
	mongoURI := os.Getenv("MONGOURI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	fmt.Println(mongoURI)
	fmt.Println(dbName)
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
	dbName := os.Getenv("DB_NAME")
	return client.Database(dbName).Collection(collectionName)
}

func Register(user User) (string, error) {
	collection := GetDBCollection("Users")

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(string); ok {
		return oid, nil
	}

	return "Registered Successfully", fmt.Errorf("registered successfully")
}

func GetUserByEmail(email string) (User, error) {
	var user User
	collection := GetDBCollection("Users")

	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("User not found")
		}
		return User{}, err
	}

	return user, nil
}
