package models

import "gorm.io/gorm"

// Artist represents an aritst entity in the database
type Artist struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Songs []Song `gorm:"foreignKey:ArtistID"`
}
