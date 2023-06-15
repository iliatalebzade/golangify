package main

import (
	"encoding/json"
	"log"
	"net/http"
	"nice_stream/config"
	"nice_stream/controllers"
	"nice_stream/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Load the application configuration
	err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure the storage paths
	err = utils.ConfigureStorage()
	if err != nil {
		log.Fatalf("Failed to configure storage: %v", err)
	}

	// Create a new router using Gorilla mux
	router := mux.NewRouter()

	// Initialize the Database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Initialize the song controller
	songController := controllers.NewSongController(db)

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
