package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nice_stream/config"
	"nice_stream/models"
	"nice_stream/utils"
	"path/filepath"

	"gorm.io/gorm"
)

// SongController handles HTTP requests related to songs.
type SongController struct {
	db *gorm.DB
}

// NewSongController creates a new instance of SongController.
func NewSongController(db *gorm.DB) *SongController {
	return &SongController{db: db}
}

// CreateSong handles the creation of a new song.
func (c *SongController) CreateSong(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Limit: 32MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrive the form fields and files
	name := r.FormValue("name")
	artistName := r.FormValue("artist")

	cover, coverHeader, err := r.FormFile("cover")
	if err != nil {
		http.Error(w, "Failed to retrieve cover image", http.StatusBadRequest)
		return
	}
	defer cover.Close()

	audio, audioHeader, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Failed to retrieve audio file", http.StatusBadRequest)
		return
	}
	defer audio.Close()

	// Create a new artist record or retrieve an existing one
	artist := models.Artist{}
	err = c.db.FirstOrCreate(&artist, models.Artist{Name: artistName}).Error
	if err != nil {
		http.Error(w, "Failed to create or retrieve artist record", http.StatusInternalServerError)
		return
	}

	// Generate unique filenames for the cover image and audio file
	coverFileName := fmt.Sprintf("%d_%s", artist.ID, coverHeader.Filename)
	audioFileName := fmt.Sprintf("%d_%s", artist.ID, audioHeader.Filename)

	// Save the cover image to the storage
	coverPath := filepath.Join(config.Config.StorageDir, "covers", coverFileName)
	err = utils.SaveFile(cover, coverPath)
	if err != nil {
		http.Error(w, "Failed to save cover image", http.StatusInternalServerError)
		return
	}

	// Save the audio file to the storage
	audioPath := filepath.Join(config.Config.StorageDir, "audio", audioFileName)
	err = utils.SaveFile(audio, audioPath)
	if err != nil {
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	// Create a new song record
	song := models.Song{
		Name:     name,
		Artist:   artist,
		CoverURL: "/covers/" + coverFileName,
		AudioURL: "/audio/" + audioFileName,
	}

	err = c.db.Create(&song).Error
	if err != nil {
		http.Error(w, "Failed to create song record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// GetSongs returns a list of songs with their details.
func (c *SongController) GetSongs(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song
	err := c.db.Preload("Artist").Find(&songs).Error
	if err != nil {
		http.Error(w, "Failed to retrieve songs", http.StatusInternalServerError)
	}

	type SongResponse struct {
		Name     string `json:"name"`
		Artist   string `json:"artist"`
		CoverURL string `json:"coverUrl"`
		AudioURL string `json:"audioUrl"`
	}

	var songResponses []SongResponse
	for _, song := range songs {
		songResponses = append(songResponses, SongResponse{
			Name:     song.Name,
			Artist:   song.Artist.Name,
			CoverURL: song.CoverURL,
			AudioURL: song.AudioURL,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songResponses)
}
