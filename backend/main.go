package main

import (
	"encoding/json"
	"log"
	"net/http"
	"nice_stream/config"
	"nice_stream/controllers"
	"nice_stream/middlewares"
	"nice_stream/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load the application configuration
	err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err.Error())
	}

	// Load the .env file to get the secret_key used in encoding and decoding tokens
	err = godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %s", err.Error())
	}

	// Configure the storage paths
	err = utils.ConfigureStorage()
	if err != nil {
		log.Fatalf("Failed to configure storage: %s", err.Error())
	}

	// Create a new router using Gorilla mux
	router := mux.NewRouter()

	// Initialize the Database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err.Error())
	}

	// Initialize the song controller
	songController := controllers.NewSongController(db)

	// Initiate the middlewares
	authorizationMiddleware := middlewares.NewAuthorize(db)

	// Define the /songs routes with authorization middleware
	songsRouter := router.PathPrefix("/songs").Subrouter()
	songsRouter.Use(authorizationMiddleware.CheckToken)
	songsRouter.HandleFunc("", songController.CreateSong).Methods("POST")
	songsRouter.HandleFunc("", songController.GetSongs).Methods("GET")

	// Initialize the user controller
	userController := controllers.NewUserController(db)

	// Define the user routes
	router.HandleFunc("/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/login", userController.LoginUser).Methods("POST")
	router.HandleFunc("/logout", userController.LogoutUser).Methods("DELETE")
	router.HandleFunc("/verify_account/{token}", userController.VerifyUser).Methods("GET")

	// Define the /songs routes
	router.HandleFunc("/songs", songController.CreateSong).Methods("POST")
	router.HandleFunc("/songs", songController.GetSongs).Methods("GET")

	// Define the root route
	router.HandleFunc("/", rootHandler).Methods("GET")

	// Serve the application
	log.Println("Server started on http://193.163.200:8000")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}

type RootResponse struct {
	Message string `json:"message"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := RootResponse{
		Message: "Welcome to the GoPotify app",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
