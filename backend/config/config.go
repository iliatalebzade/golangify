package config

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Configuration represents the application configuration.
type Configuration struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`

	SmppServer struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smppServer"`

	StorageDir string `json:"storageDir"`
	CoversDir  string `json:"coversDir"`
	AudioDir   string `json:"audioDir"`
}

var (
	// Config stores the application configuration.
	Config Configuration
)

// LoadConfig loads the application configuration from the specified file.
func LoadConfig(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig returns a reference to the application configuration.
func GetConfig() (*Configuration, error) {
	// Check if the Config variable is empty
	if Config.Database.Host == "" {
		// If empty, load the configuration from the file
		err := LoadConfig("config/config.json")
		if err != nil {
			return nil, err
		}
	}

	// Return the instance of the Configuration
	return &Config, nil
}

// InitDB initializes the database connection using the configuration.
func InitDB() (*gorm.DB, error) {
	dsn := Config.Database.User + ":" + Config.Database.Password + "@tcp(" +
		Config.Database.Host + ":" + strconv.Itoa(Config.Database.Port) + ")/" + Config.Database.Name +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run the migrations
	err = migrateModels(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
