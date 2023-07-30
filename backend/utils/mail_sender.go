package utils

import (
	"bytes"
	"nice_stream/config"
	"nice_stream/models"
	"os"
	"text/template"

	"gopkg.in/gomail.v2"
)

const (
	from    string = "account_verifier@gopotify.io"
	subject string = "Account Confirmation"
)

type TemplateData struct {
	Token string
	Email string
}

// SendConfirmationMail sends an email to the recipient user containing the token, username, and email for confirmation.
func SendConfirmationMail(recipient *models.User, verificationToken string) error {
	// Get the configs from the json file
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Replace these values with your actual email configuration.
	password := config.SmppServer.Password
	smtpHost := config.SmppServer.Host
	smtpPort := config.SmppServer.Port

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipient.Email)
	m.SetHeader("Subject", subject)

	// Read the HTML template file
	templateFilePath := "./utils/templates/mails/account_verification.html"
	templateContent, err := os.ReadFile(templateFilePath)
	if err != nil {
		return err
	}

	// Parse the HTML template
	tmpl, err := template.New("account_verification").Parse(string(templateContent))
	if err != nil {
		return err
	}

	// Prepare the data to be used in the template
	data := TemplateData{
		Token: verificationToken,
		Email: recipient.Email,
	}

	// Render the template into a buffer
	var emailBodyBuf bytes.Buffer
	err = tmpl.Execute(&emailBodyBuf, data)
	if err != nil {
		return err
	}

	// Set the email body as HTML
	m.SetBody("text/html", emailBodyBuf.String())

	// Create the SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email using the dialer
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
