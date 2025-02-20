package notification

import (
	"context"
	"loan_service/models"
)

// Service defines the interface for notification operations
type Service interface {
	// SendAgreementLetter sends agreement letter to investor
	SendAgreementLetter(ctx context.Context, loan *models.Loan, investor *models.Investor, agreementURL string) error

	// SendInvestmentConfirmation sends investment confirmation to investor
	SendInvestmentConfirmation(ctx context.Context, investment *models.Investment) error

	// SendDisbursementNotice sends disbursement notice to all investors
	SendDisbursementNotice(ctx context.Context, loan *models.Loan) error
}

// EmailData represents the data needed for email templates
type EmailData struct {
	RecipientName string
	LoanID        int64
	LoanAmount    float64
	InvestAmount  float64
	ROI           float64
	DocumentURL   string
}
