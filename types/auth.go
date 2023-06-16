package types

import "github.com/golang-jwt/jwt"

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string  `json:"message"`
	Token   *string `json:"token"`
}

type AuthenticateResponse struct {
	Message string `json:"message"`
}
