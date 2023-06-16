package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"nice_stream/types"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type contextKey string

const (
	ContextKeyUsername contextKey = "username"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Search the request headers for an `Authorization` header
		// And if successful get the value
		tokenHeader := r.Header.Get("Authorization")

		// A message to use whenever the authorize process was unsuccessful (the use was unauthorized)
		unauthorized_response := types.AuthenticateResponse{
			Message: "You're not authorized to visit this page",
		}

		// Split the Authorization header value to separate the token type and token
		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(unauthorized_response)
			return
		}

		// Select the part of the string that contains the token and the token only
		tokenString := parts[1]

		// Check if the token part exists and if not,
		// Return with a StatusUnauthorized, informing the client of their lack of access to the requested page
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(unauthorized_response)
			return
		}

		// Load the .env file to get the secret_key used in encoding and decoding tokens
		err := godotenv.Load()
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
			json.NewEncoder(w).Encode(unauthorized_response)
			return
		}

		// Parse the content of the token into a Claims struct
		// If the token was not valid or encountered any other errors
		// Return with a StatusUnauthorized, informing the client of their lack of access to the requested page
		claims, ok := token.Claims.(*types.Claims)
		if !ok || !token.Valid {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(unauthorized_response)
			return
		}

		// Pass the username to the next handler as a context value
		r = r.WithContext(context.WithValue(r.Context(), ContextKeyUsername, claims.Email))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
