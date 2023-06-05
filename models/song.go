package models

import "gorm.io/gorm"

// Song represents a song entity in the database
type Song struct {
	gorm.Model
	Name     string
	ArtistID uint
	Artist   Artist `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CoverURL string
	AudioURL string
}
