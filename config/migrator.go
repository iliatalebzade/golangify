package config

import (
	"nice_stream/models"

	"gorm.io/gorm"
)

func migrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&models.Song{}, &models.Artist{}, &models.User{})
}
