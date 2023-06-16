package utils

import (
	"log"
	"nice_stream/types"
	"os"
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

	// Replace "your-secret-key" with your actual secret key
	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
