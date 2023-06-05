package utils

import (
	"nice_stream/config"
	"os"
	"path/filepath"
)

// ConfigureStorage initializes the storage paths and creates the required directories.
func ConfigureStorage() error {
	config.Config.StorageDir = filepath.Join(".", "storage")
	config.Config.CoversDir = filepath.Join(config.Config.StorageDir, "covers")
	config.Config.AudioDir = filepath.Join(config.Config.StorageDir, "audio")

	// Create the storage directories if they don't exist
	err := os.MkdirAll(config.Config.CoversDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(config.Config.AudioDir, 0755)
	if err != nil {
		return err
	}

	return nil
}
