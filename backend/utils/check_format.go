package utils

import "net/mail"

// Checks if the string containing an email address is valid or not
func IsValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
