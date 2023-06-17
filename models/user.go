package models

import "gorm.io/gorm"

// User represents a user entity in the database
type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(255);unique"`
	Password string
}

// BlacklistedToken represents a blacklisted token stored in the database.
type BlacklistedToken struct {
	gorm.Model
	Token string `gorm:"type:varchar(255);uniqueIndex"`
}
