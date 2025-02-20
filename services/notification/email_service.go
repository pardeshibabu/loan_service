package notification

import (
	"context"
	"fmt"
	"time"

	"loan_service/models"

	"github.com/gobuffalo/pop/v6"
)

type emailService struct {
	db     *pop.Connection
	mailer *Mailer
}

// NewEmailService creates a new email notification service
func NewEmailService(db *pop.Connection) (Service, error) {
	config, err := NewEmailConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create email config: %w", err)
	}

	return &emailService{
		db:     db,
		mailer: NewMailer(config),
	}, nil
}

func (s *emailService) SendAgreementLetter(ctx context.Context, loan *models.Loan, investor *models.Investor, agreementURL string) error {
	notification := &models.Notification{
		LoanID:       loan.ID,
		InvestorID:   &investor.ID,
		Type:         models.NotificationTypeAgreementLetter,
		Status:       models.NotificationStatusPending,
		EmailContent: fmt.Sprintf("Agreement letter for loan #%d", loan.ID),
	}

	if err := s.db.Create(notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send email
	data := EmailData{
		RecipientName: fmt.Sprintf("%s %s", investor.FirstName, investor.LastName),
		LoanAmount:    loan.PrincipalAmount,
		DocumentURL:   agreementURL,
	}

	err := s.mailer.SendEmail(
		investor.Email,
		"Loan Agreement Letter",
		"agreement_letter",
		data,
	)

	now := time.Now()
	if err != nil {
		notification.Status = models.NotificationStatusFailed
		notification.EmailContent = err.Error()
	} else {
		notification.Status = models.NotificationStatusSent
		notification.SentAt = &now
	}

	return s.db.Update(notification)
}

func (s *emailService) SendInvestmentConfirmation(ctx context.Context, investment *models.Investment) error {
	// Get loan and investor details
	loan := &models.Loan{}
	if err := s.db.Find(loan, investment.LoanID); err != nil {
		return fmt.Errorf("failed to find loan: %w", err)
	}

	investor := &models.Investor{}
	if err := s.db.Find(investor, investment.InvestorID); err != nil {
		return fmt.Errorf("failed to find investor: %w", err)
	}

	notification := &models.Notification{
		LoanID:       investment.LoanID,
		InvestorID:   &investment.InvestorID,
		Type:         models.NotificationTypeInvestmentConfirmation,
		Status:       models.NotificationStatusPending,
		EmailContent: fmt.Sprintf("Investment confirmation for loan #%d", investment.LoanID),
	}

	if err := s.db.Create(notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send email
	data := EmailData{
		RecipientName: fmt.Sprintf("%s %s", investor.FirstName, investor.LastName),
		LoanID:        loan.ID,
		LoanAmount:    loan.PrincipalAmount,
		InvestAmount:  investment.Amount,
		ROI:           loan.ROI,
	}

	err := s.mailer.SendEmail(
		investor.Email,
		"Investment Confirmation",
		"investment_confirmation",
		data,
	)

	now := time.Now()
	if err != nil {
		notification.Status = models.NotificationStatusFailed
		notification.EmailContent = err.Error()
	} else {
		notification.Status = models.NotificationStatusSent
		notification.SentAt = &now
	}

	return s.db.Update(notification)
}

func (s *emailService) SendDisbursementNotice(ctx context.Context, loan *models.Loan) error {
	// Get all investments with investors
	var investments []models.Investment
	if err := s.db.Where("loan_id = ?", loan.ID).All(&investments); err != nil {
		return fmt.Errorf("failed to get investments: %w", err)
	}

	for _, investment := range investments {
		investor := &models.Investor{}
		if err := s.db.Find(investor, investment.InvestorID); err != nil {
			return fmt.Errorf("failed to find investor: %w", err)
		}

		notification := &models.Notification{
			LoanID:       loan.ID,
			InvestorID:   &investment.InvestorID,
			Type:         models.NotificationTypeDisbursementNotice,
			Status:       models.NotificationStatusPending,
			EmailContent: fmt.Sprintf("Loan #%d has been disbursed", loan.ID),
		}

		if err := s.db.Create(notification); err != nil {
			return fmt.Errorf("failed to create notification: %w", err)
		}

		// Send email
		data := EmailData{
			RecipientName: fmt.Sprintf("%s %s", investor.FirstName, investor.LastName),
			LoanID:        loan.ID,
			LoanAmount:    loan.PrincipalAmount,
			InvestAmount:  investment.Amount,
			DocumentURL:   *loan.SignedAgreementURL,
		}

		err := s.mailer.SendEmail(
			investor.Email,
			"Loan Disbursement Notice",
			"disbursement_notice",
			data,
		)

		now := time.Now()
		if err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.EmailContent = err.Error()
		} else {
			notification.Status = models.NotificationStatusSent
			notification.SentAt = &now
		}

		if err := s.db.Update(notification); err != nil {
			return fmt.Errorf("failed to update notification: %w", err)
		}
	}

	return nil
}
