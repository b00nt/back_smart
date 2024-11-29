package services

import (
	"back/internal/models"
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendEmail(feedback models.Feedback) error {
	from := "b00nt1377@gmail.com"     // Replace with your email address
	to := "matkovsky88@protonmail.ch" // Replace with the recipient's email address
	subject := "New Feedback Submission"
	body := fmt.Sprintf("New feedback received:\n\nName: %s\nEmail: %s\nMessage: %s", feedback.Name, feedback.Telephone, feedback.City)
	psswd := "Lm=%cy7{ARH~NHY6(s_h@q"

	// Set up the email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// Configure the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, from, psswd) // Use app password if 2FA is on

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
