package middleware

import (
	"complaints/cmd/api/models"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var jwtKeyEncoded = os.Getenv("JWTKEY")

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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
			fmt.Println("Claims")
			return
		}

		email := claims["iss"].(string)

		// Get matric number
		matricNo, err := models.GetMatricNo(email)
		if err != nil {
			fmt.Println("Unavle to get matric number")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//store matric number in the request context
		ctx := context.WithValue(r.Context(), "matricNo", matricNo)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
