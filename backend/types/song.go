package types

// A struct representing a single SongItem
type SongItem struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	CoverURL string `json:"coverUrl"`
	AudioURL string `json:"audioUrl"`
}

// An array of SongItem objects
type SongsList []SongItem

// A list of songs returned in response to GetSongs
type GetSongsResponse struct {
	// All songs in the system
	Songs []SongItem `json:"songs"`
}
