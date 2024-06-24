package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	Password  string             `json:"password" bson:"password"`
}

type LoginCredentials struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password" bson:"password"`
}

type Student struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	MatricNo  string             `json:"matric_no" bson:"matric_no"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Courses   []string           `json:"courses,omitempty" bson:"courses,omitempty"`
	Program   string
}

type Lecturer struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	StaffID      string             `json:"staff_id" bson:"staff_id"`
	FirstName    string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName     string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	CoursesTaken []string           `json:"courses_taken,omitempty" bson:"courses_taken,omitempty"`
}

type Course struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CourseCode       string             `json:"course_code" bson:"course_code"`
	CourseName       string             `json:"course_name,omitempty" bson:"course_name,omitempty"`
	Semester         string             `json:"semester,omitempty" bson:"semester,omitempty"`
	StudentsEnrolled []string           `json:"students_enrolled,omitempty" bson:"students_enrolled,omitempty"`
	Lecturers        []string           `json:"lecturers,omitempty" bson:"lecturers,omitempty"`
}

type Complaint struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	RequestingStudent  string             `json:"requesting_student,omitempty" bson:"requesting_student,omitempty"`
	RequestDetails     string             `json:"request_details,omitempty" bson:"request_details,omitempty"`
	StudentProof       string             `json:"student_proof,omitempty" bson:"student_proof,omitempty"`
	LecturerProof      string             `json:"lecturer_proof,omitempty" bson:"lecturer_proof,omitempty"`
	TestScore          int                `json:"test_score,omitempty" bson:"test_score,omitempty"`
	CourseConcerned    string             `json:"course_concerned,omitempty" bson:"course_concerned,omitempty"`
	RespondingLecturer string             `json:"responding_lecturer,omitempty" bson:"responding_lecturer,omitempty"`
	Status             string             `json:"status,omitempty" bson:"status,omitempty"`
	Reason             string             `json:"reason,omitempty" bson:"reason,omitempty"`
	CreatedAt          time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Senate struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password"`
}

type CSIS struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password"`
}
