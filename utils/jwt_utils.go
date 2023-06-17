package utils

import (
	"errors"
	"log"
	"net/http"
	"nice_stream/types"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Generates json-web-tokens and returns it
func GenerateToken(email string) (string, error) {
	claims := &types.Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 72).Unix(),
		},
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return "", nil
	}

	secret_key := os.Getenv("TOKEN_SALT")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header.
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("invalid token format")
	}

	return tokenParts[1], nil
}
