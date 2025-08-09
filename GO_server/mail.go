package main

import (
	"fmt"
	"net/smtp"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

func SendMail(to []string, subject, body string) error {
	// Use environment variables or defaults
	smtpHost := getEnv("SMTP_HOST", "smtp.example.com")
	smtpPort := getEnvAsInt("SMTP_PORT", 587)
	smtpUser := getEnv("SMTP_USER", "user")
	smtpPass := getEnv("SMTP_PASS", "pass")

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	return d.DialAndSend(m)
}

// simple alternative using net/smtp (not used but handy)
func sendSimpleSMTP(to []string, subject, body string) error {
	smtpHost := getEnv("SMTP_HOST", "smtp.example.com")
	smtpPort := getEnv("SMTP_PORT", "587")
	smtpUser := getEnv("SMTP_USER", "user")
	smtpPass := getEnv("SMTP_PASS", "pass")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	msg := "From: " + smtpUser + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	return smtp.SendMail(addr, auth, smtpUser, to, []byte(msg))
}
