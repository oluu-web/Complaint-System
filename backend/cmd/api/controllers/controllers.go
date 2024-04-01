package controllers

import (
	"complaints/cmd/api/models"
	"complaints/cmd/api/utilities"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

var jwtKeyEncoded = os.Getenv("JWTKEY")

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//store in the database
	user.Password = string(hashedPassword)
	_, err = models.Register(user)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, user, "user")
}

func Login(w http.ResponseWriter, r *http.Request) {

	var credentials models.LoginCredentials
	type jsonResp struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
		UserID  string `json:"user_id"`
		Token   string `json:"token"`
		Role    string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bad := jsonResp{
		OK:      false,
		Message: "Invalid username or password",
		UserID:  "null",
	}

	//find user by ID
	user, err := models.GetUserByUserID(credentials.Username)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, bad, "response")
		return
	}

	//compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//Generste JWT token
	jwtKey, err := base64.URLEncoding.DecodeString(jwtKeyEncoded)
	fmt.Println("JWT KEY ENCODED, LINE 98: ", jwtKeyEncoded)
	if err != nil {
		http.Error(w, "Error decoding JWT key", http.StatusInternalServerError)
		fmt.Println("Error decoding JWT key:", err)
		return
	}
	type CustomClaims struct {
		jwt.StandardClaims
		Role string `json:"role,omitempty"`
	}
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Subject:   user.ID.Hex(),
			Issuer:    credentials.Username,
		},
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("line 114: ", jwtKey)
		return
	}
	ok := jsonResp{
		OK:      true,
		Message: "Login successful",
		UserID:  credentials.Username,
		Token:   tokenString,
		Role:    user.Role,
	}

	//send token in response
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
	utilities.WriteJSON(w, http.StatusOK, ok, "response")
	fmt.Println("Logged in successfully!")
}

func NewComplaint(w http.ResponseWriter, r *http.Request) {
	var complaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&complaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	studentId, ok := r.Context().Value("userID").(string)
	if !ok {
		utilities.ErrorJSON(w, errors.New("unable to get student ID from context"))
		return
	}
	complaint.RequestingStudent = studentId
	course, err := models.GetCourseByCourseCode(string(complaint.CourseConcerned))
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
	if len(course.Lecturers) == 0 {
		utilities.ErrorJSON(w, errors.New("no lecturers found for the course"))
		return
	}
	respondingLecturer := course.Lecturers[rand.Intn(len(course.Lecturers))]

	complaint.RespondingLecturer = respondingLecturer
	complaint.Status = "Pending"
	complaint.CreatedAt = time.Now()
	complaint.UpdatedAt = time.Now()

	_, err = models.CreateNewComplaint(complaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}
}

func ExtractEmailFromRequest(r *http.Request) string {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return ""
	}
	return requestData.Email
}

func GetComplaintByObjectID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	complaint, err := models.GetComplaintByObjectId(objID)
	if err != nil {
		fmt.Println("Unable to get complaint")
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaint, "complaint")
}

func GetComplaintsByStaffID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	complaints, err := models.GetComplaintsByStaffId(id)
	if err != nil {
		fmt.Println("Unable to get complaint")
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func GetComplaintsByStudentID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	complaints, err := models.GetComplaintsByStudentId(id)
	if err != nil {
		fmt.Println("Unable to get complaint")
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, complaints, "complaints")
}

func ChangeComplaintStatusByLecturer(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By Lecturer")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func ChangeComplaintStatusByAdvisor(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By Course Advisor")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func ChangeComplaintStatusByHOD(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By HOD")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func ChangeComplaintStatusBySenate(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var updatedComplaint models.Complaint
	err := json.NewDecoder(r.Body).Decode(&updatedComplaint)
	if err != nil {
		utilities.ErrorJSON(w, err)
		return
	}

	err = models.ChangeComplaintStatus(id, "Approved By Senate")
	if err != nil {
		utilities.ErrorJSON(w, err)
	}

	utilities.WriteJSON(w, http.StatusOK, "Status Updated Successfully", "Success")
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	student, err := models.GetStudentById(id)
	if err != nil {
		fmt.Println("Unable to get student")
		utilities.ErrorJSON(w, err)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, student, "student")
}

func GetCoursesByStudentID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	student, err := models.GetStudentById(id)
	if err != nil {
		fmt.Println("Unable to get student")
		utilities.ErrorJSON(w, err)
		return
	}

	courses := student.Courses

	utilities.WriteJSON(w, http.StatusOK, courses, "courses")
}
