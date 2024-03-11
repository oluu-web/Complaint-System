package controllers

import (
	"complaints/cmd/api/models"
	"complaints/cmd/api/utilities"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//find user by email
	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
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
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//send token in response
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
	fmt.Println("Logged in successfully!")
}
