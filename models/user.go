package models

import "gorm.io/gorm"

// User represents a user entity in the database
type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(255);unique"`
	Password string
}
