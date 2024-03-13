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
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matricNo, err := models.GetMatricNo(credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ok := jsonResp{
		OK:      true,
		Message: "Login successful",
		UserID:  matricNo,
	}

	bad := jsonResp{
		OK:      false,
		Message: "Invalid email or password",
		UserID:  "null",
	}

	//find user by email
	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, bad, "response")
		return
	}

	//compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	//Generste JWT token
	jwtKey, err := base64.URLEncoding.DecodeString(jwtKeyEncoded)
	if err != nil {
		fmt.Println("Error decoding JWT key:", err)
		return
	}
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   user.ID.Hex(),
		Issuer:    credentials.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	fmt.Println("Line 115, JWT key", jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(tokenString)

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

	studentId, ok := r.Context().Value("matricNo").(string)
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
