package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

// var jwtKeyEncoded = os.Getenv("JWTKEY")
var jwtKeyEncoded = "GQFUUfN75vQdsYvzJBmXEQvICiX9HU8HfHrPkNJfRq0="

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtKey, err := base64.URLEncoding.DecodeString(jwtKeyEncoded)
		if err != nil {
			fmt.Println("Error decoding JWT key:", err)
			return
		}

		//Get token from request header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			fmt.Println("Empty auth")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//parse and validate token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexected signing method")
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {

			fmt.Println("Token not valid: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//extract email from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims["iss"].(string)

		//store matric number in the request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
