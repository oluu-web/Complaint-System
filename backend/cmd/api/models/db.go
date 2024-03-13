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

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	oid := user.ID.Hex()

	return oid, nil
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

	matNo, err := GetMatricNo(email)
	if err != nil {
		fmt.Println("Unable to get Matric Number")
	}
	user.MatricNo = matNo

	return user, nil
}

func CreateNewComplaint(complaint Complaint) (string, error) {
	collection := GetDBCollection("Complaints")

	_, err := collection.InsertOne(context.Background(), complaint)
	if err != nil {
		return "", fmt.Errorf("failed to insert complaint: %w", err)
	}

	oid := complaint.ID.Hex()

	return oid, nil
}

func GetCourseByCourseCode(courseCode string) (Course, error) {
	collection := GetDBCollection("Courses")

	filter := bson.M{
		"course_code": courseCode,
	}

	var course Course
	err := collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Course{}, fmt.Errorf("Course not found")
		}
		return Course{}, err
	}
	return course, nil
}

func GetMatricNo(email string) (string, error) {
	collection := GetDBCollection("Students")

	filter := bson.M{
		"email": email,
	}

	var result struct {
		MatricNo string `bson:"matric_no"`
	}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.MatricNo, nil
}
