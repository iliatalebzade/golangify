package main

import (
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

	// Define the routes
	router.HandleFunc("/songs", songController.CreateSong).Methods("POST")
	router.HandleFunc("/songs", songController.GetSongs).Methods("GET")

	// Serve the application
	log.Println("Server started on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
