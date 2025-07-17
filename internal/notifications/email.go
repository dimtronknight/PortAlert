package notifications

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailNotifier struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	ToEmail      string
}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		ToEmail:      os.Getenv("TO_EMAIL"),
	}
}

func (e *EmailNotifier) SendEmail(subject, body string) error {
	if e.SMTPHost == "" || e.SMTPUsername == "" || e.SMTPPassword == "" {
		return fmt.Errorf("SMTP credentials are required")
	}

	// Set default port if not provided
	if e.SMTPPort == "" {
		e.SMTPPort = "587"
	}

	// Use username as from email if not provided
	if e.FromEmail == "" {
		e.FromEmail = e.SMTPUsername
	}

	// Use from email as to email if not provided (for testing)
	if e.ToEmail == "" {
		e.ToEmail = e.FromEmail
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.FromEmail, e.ToEmail, subject, body)

	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)

	err := smtp.SendMail(
		e.SMTPHost+":"+e.SMTPPort,
		auth,
		e.FromEmail,
		[]string{e.ToEmail},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailNotifier) SendInvestmentUpdate(totalBGN, stocksBGN, cryptoBGN, bullionBGN float64) error {
	subject := "💰 Daily Investment Update"
	body := fmt.Sprintf(
		"Good morning! Here's your daily portfolio update:\n\n"+
			"📊 Portfolio Breakdown:\n"+
			"🏦 Trading212: %.2f BGN\n"+
			"₿ Crypto: %.2f BGN\n"+
			"🥇 Bullion: %.2f BGN\n\n"+
			"💎 Total Investment Worth: %.2f BGN\n\n"+
			"📈 Have a great day!\n\n"+
			"---\n"+
			"This is an automated daily update from your Investment Tracker",
		stocksBGN, cryptoBGN, bullionBGN, totalBGN,
	)

	return e.SendEmail(subject, body)
}