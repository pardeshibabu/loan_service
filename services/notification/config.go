package notification

import (
	"fmt"
	"os"
	"strconv"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
}

// NewEmailConfig creates email config from environment variables
func NewEmailConfig() (*EmailConfig, error) {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}

	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	return &EmailConfig{
		Host:     host,
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: os.Getenv("EMAIL_FROM_NAME"),
		FromAddr: os.Getenv("EMAIL_FROM_ADDRESS"),
	}, nil
}
