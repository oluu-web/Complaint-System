package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	Password  string             `json:"password" bson:"password"`
}

type Student struct {
	MatricNo  string   `json:"matric_no" bson:"matric_no"`
	FirstName string   `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string   `json:"email,omitempty" bson:"email,omitempty"`
	Courses   []Course `json:"courses,omitempty" bson:"courses,omitempty"`
	Program   string
}

type Lecturer struct {
	StaffID      string   `json:"staff_id" bson:"staff_id"`
	FirstName    string   `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName     string   `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email        string   `json:"email,omitempty" bson:"email,omitempty"`
	CoursesTaken []Course `json:"courses_taken,omitempty" bson:"courses_taken,omitempty"`
}

type Course struct {
	CourseCode       string     `json:"course_code" bson:"course_code"`
	CourseName       string     `json:"course_name,omitempty" bson:"course_name,omitempty"`
	Semester         string     `json:"semester,omitempty" bson:"semester,omitempty"`
	StudentsEnrolled []Student  `json:"students_enrolled,omitempty" bson:"students_enrolled,omitempty"`
	Lecturers        []Lecturer `json:"lecturers,omitempty" bson:"lecturers,omitempty"`
}

type Request struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
	RequestingStudent  string             `json:"requesting_student,omitempty" bson:"requesting_student,omitempty"`
	RequestDetails     string             `json:"request_details,omitempty" bson:"request_details,omitempty"`
	FilePath           string             `json:"file_path,omitempty" bson:"file_path,omitempty"`
	TestScore          int                `json:"test_score,omitempty" bson:"test_score,omitempty"`
	CourseConcerned    string             `json:"course_concerned,omitempty" bson:"course_concerned,omitempty"`
	RespondingLecturer string             `json:"responding_lecturer,omitempty" bson:"responding_lecturer,omitempty"`
	Status             string             `json:"status,omitempty" bson:"status,omitempty"`
	TimeSent           time.Time          `json:"time_sent,omitempty" bson:"time_sent,omitempty"`
}

type Senate struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password"`
}

type CSIS struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password"`
}