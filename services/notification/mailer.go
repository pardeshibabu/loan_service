package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

// Mailer handles email sending
type Mailer struct {
	config *EmailConfig
	auth   smtp.Auth
}

// NewMailer creates a new mailer instance
func NewMailer(config *EmailConfig) *Mailer {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &Mailer{
		config: config,
		auth:   auth,
	}
}

// SendEmail sends an HTML email
func (m *Mailer) SendEmail(to string, subject string, templateName string, data interface{}) error {
	// Load template
	tmpl, err := template.ParseFiles(filepath.Join("templates", "mail", templateName+".html"))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Compose email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + body.String())

	// Send email
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	if err := smtp.SendMail(addr, m.auth, m.config.FromAddr, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
