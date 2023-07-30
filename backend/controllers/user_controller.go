package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nice_stream/models"
	"nice_stream/types"
	"nice_stream/utils"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	tokenLength int = 150
)

// A struct of the user_controller containin the database instance
type UserController struct {
	db *gorm.DB
}

// NewUserController creates a new instance of UserController.
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {

	// Receive and format and the request body object
	var req types.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the necessary parameters are indeed sent by the client
	if req.Email == "" && req.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.RegisterResponse{Message: "Please provide both email and password"})
		return
	} else if req.Email == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.RegisterResponse{Message: "Please provide an email"})
		return
	} else if req.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.RegisterResponse{Message: "Please provide a password"})
		return
	}

	// Check if the provided email address is a valid one
	if utils.IsValidEmail(req.Email) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.RegisterResponse{Message: "Please provide a valid email"})
		return
	}

	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.RegisterResponse{
			Message: "Faced an error during the registeration of your account, please try again later.",
		})
		return
	}

	verificationToken, _, err := utils.GenerateRandomToken()
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.RegisterResponse{
			Message: "Faced an error during the registeration of your account, please try again later.",
		})
		return
	}

	// Create the new User object
	user := models.User{
		Email:         req.Email,
		Password:      string(hashedPassword),
		VerifiedToken: verificationToken,
	}

	// Attempt to save the User object to the database
	err = uc.db.Create(&user).Error
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), fmt.Sprintf("Duplicate entry '%s' for key 'users.email'", req.Email)) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(types.RegisterResponse{
				Message: "This email is already in use, please sign in, change password or use another one.",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.RegisterResponse{
			Message: "Faced an error during the registeration of your account, please try again later.",
		})
		return
	}

	// Send an email to the provided email to verify their account
	err = utils.SendConfirmationMail(&user, verificationToken)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.RegisterResponse{
			Message: "Faced an error during the registeration of your account, please try again later.",
		})

		return
	}

	// Inform the Client of the registrations success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.RegisterResponse{Message: "Registration successful! please verify your account to be able to use it."})
}

func (uc *UserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Extract the parameter variables out of the request into a hash
	vars := mux.Vars(r)

	// From the request params hash, select the token `slug`
	token := vars["token"]

	// Get the salt string used to encode the token
	saltString := os.Getenv("EMAIL_VERIFICATION_SALT")

	// Decode the token
	decodedTokenWithSalt, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		// Handle decoding error
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.VerifyUserResponse{
			Message: "Token was invalid",
		})
		return
	}

	// Remove the salt from the decoded token
	decodedTokenBytes, err := base64.URLEncoding.DecodeString(string(decodedTokenWithSalt[len(saltString):]))
	if err != nil {
		// Handle decoding error
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.VerifyUserResponse{
			Message: "Token was invalid",
		})
		return
	}

	// Generate the expected token
	_, tokenBytes, _ := utils.GenerateRandomToken()
	expectedTokenBytes := []byte(saltString + base64.URLEncoding.EncodeToString(tokenBytes))
	expectedToken := base64.URLEncoding.EncodeToString(expectedTokenBytes)

	// Compare the decoded token with the expected token
	if string(decodedTokenBytes) != expectedToken {
		// Handle invalid token
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.VerifyUserResponse{
			Message: "Token was invalid",
		})
		return
	}

	// Token is valid, perform further work
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(types.VerifyUserResponse{
		Message: "Token was valid",
	})
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {

	// Receive and format and the request body object
	var req types.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the necessary parameters are indeed sent by the client
	if req.Email == "" && req.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.LoginResponse{
			Message: "Please provide both email and password",
			Token:   nil,
		})
		return
	} else if req.Email == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.LoginResponse{
			Message: "Please provide an email",
			Token:   nil,
		})
		return
	} else if req.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.LoginResponse{
			Message: "Please provide a password",
			Token:   nil,
		})
		return
	}

	// Check if the provided email address is a valid one
	if utils.IsValidEmail(req.Email) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.RegisterResponse{
			Message: "Please provide a valid email",
		})
		return
	}

	// Create a user instance to search the database for the correct instance
	// And if found, populate it with the data
	var user models.User
	err = uc.db.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		log.Println(err.Error())
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.LoginResponse{
				Message: "Couldn't find an email with the given combination, try siging up or change password",
				Token:   nil,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the provided password and the hash of found user instance match
	// (if the password is correct)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.LoginResponse{
			Message: "Couldn't find an email with the given combination, try siging up or change password",
			Token:   nil,
		})
		return
	}

	// Generate the JWT to send it to the client
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.LoginResponse{
			Message: "Encountered an error trying to sign you in, please try again later",
			Token:   nil,
		})
		return
	}

	// Inform the client of the success and providing them with the generated token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.LoginResponse{Message: "Login was successful.", Token: &token})
}

func (uc *UserController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the request headers
	token, err := utils.ExtractTokenFromHeader(r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.AuthenticateResponse{
			Message: "You're not authorized to visit this page",
		})
		return
	}

	// Create a new BlacklistedToken object
	blacklistedToken := models.BlacklistedToken{
		Token: token,
	}

	// Save the blacklisted token to the database
	err = uc.db.Create(&blacklistedToken).Error
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.LogoutResponse{
			Message: "Encountered an error trying to sign you out, please try again later",
		})
		return
	}

	// Respond with a successful logout message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.LogoutResponse{
		Message: "Logout successful.",
	})
}
