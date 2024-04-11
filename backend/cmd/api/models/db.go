package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	mongoURI := os.Getenv("MONGOURI")
	if mongoURI == "" {
		return fmt.Errorf("MONGOURI environment variable not set")
	}
	fmt.Println("MongoURI: ", mongoURI)
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

func GetUserByObjectID(id primitive.ObjectID) (User, error) {
	var user User
	collection := GetDBCollection("Users")

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("User not found")
		}
		return User{}, err
	}
	return user, nil
}

func GetUserByUserID(userID string) (User, error) {
	var user User
	collection := GetDBCollection("Users")

	err := collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("User not found")
		}
		return User{}, err
	}
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

func GetComplaintByObjectId(id primitive.ObjectID) (Complaint, error) {
	var complaint Complaint
	collection := GetDBCollection("Complaints")
	filter := bson.M{
		"_id": id,
	}

	err := collection.FindOne(context.Background(), filter).Decode(&complaint)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Complaint{}, fmt.Errorf("Complaint not found")
		}
		return Complaint{}, err
	}
	return complaint, nil
}

func GetComplaintsByStaffId(id string) ([]Complaint, error) {
	collection := GetDBCollection("Complaints")

	filter := bson.M{"responding_lecturer": id}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var complaints []Complaint
	for cursor.Next(context.Background()) {
		var complaint Complaint
		err := cursor.Decode(&complaint)
		if err != nil {
			return nil, err
		}
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

func GetComplaintsByStudentId(id string) ([]Complaint, error) {
	collection := GetDBCollection("Complaints")
	filter := bson.M{
		"requesting_student": id,
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var complaints []Complaint
	for cursor.Next(context.Background()) {
		var complaint Complaint
		err := cursor.Decode(&complaint)
		if err != nil {
			return nil, err
		}
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

func ChangeComplaintStatus(id string, newStatus string) error {
	collection := GetDBCollection("Complaints")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"status": newStatus,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func ChangeStatusToByHOD(id string) error {
	collection := GetDBCollection("Complaints")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"status": "Approved by HOD",
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetStudentById(userID string) (Student, error) {
	var student Student
	collection := GetDBCollection("Students")

	err := collection.FindOne(context.Background(), bson.M{"matric_no": userID}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Student{}, fmt.Errorf("Student not found")
		}
		return Student{}, err
	}
	return student, nil
}
