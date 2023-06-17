package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"nice_stream/models"
	"nice_stream/types"
	"nice_stream/utils"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type contextKey string

const (
	ContextKeyUsername contextKey = "username"
)

type Authorize struct {
	db *gorm.DB
}

func NewAuthorize(db *gorm.DB) *Authorize {
	return &Authorize{db: db}
}

func (a *Authorize) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.ExtractTokenFromHeader(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.AuthenticateResponse{
				Message: "You're not authorized to visit this page",
			})
			return
		}

		// Check if the token part exists and if not,
		// Return with a StatusUnauthorized, informing the client of their lack of access to the requested page
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.AuthenticateResponse{
				Message: "You're not authorized to visit this page",
			})
			return
		}

		// Load the .env file to get the secret_key used in encoding and decoding tokens
		err = godotenv.Load()
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Select the secret_key from the .env file
		token_salt := os.Getenv("TOKEN_SALT")

		// Check if the provided token is a valid one
		// And if so, parse it
		token, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(token *jwt.Token) (any, error) {
			return []byte(token_salt), nil
		})

		// Check if the prior part was successful and if not
		// Return with a StatusUnauthorized, informing the client of their lack of access to the requested page
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.AuthenticateResponse{
				Message: "You're not authorized to visit this page",
			})
			return
		}

		// Parse the content of the token into a Claims struct
		// If the token was not valid or encountered any other errors
		// Return with a StatusUnauthorized, informing the client of their lack of access to the requested page
		claims, ok := token.Claims.(*types.Claims)
		if !ok || !token.Valid {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.AuthenticateResponse{
				Message: "You're not authorized to visit this page",
			})
			return
		}

		// Query the blacklist tokens table to check if the token is blacklisted
		var blacklistedToken models.BlacklistedToken
		err = a.db.Where("token = ?", tokenString).First(&blacklistedToken).Error
		if err == nil {
			// Token is blacklisted
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.AuthenticateResponse{
				Message: "You're not authorized to visit this page",
			})
			return
		}

		// Pass the username to the next handler as a context value
		r = r.WithContext(context.WithValue(r.Context(), ContextKeyUsername, claims.Email))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
