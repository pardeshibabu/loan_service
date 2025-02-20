package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeAgreementLetter        NotificationType = "agreement_letter"
	NotificationTypeInvestmentConfirmation NotificationType = "investment_confirmation"
	NotificationTypeDisbursementNotice     NotificationType = "disbursement_notice"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

// Notification model
type Notification struct {
	Model
	LoanID       int64              `json:"loan_id" db:"loan_id"`
	InvestorID   *int64             `json:"investor_id" db:"investor_id"`
	Type         NotificationType   `json:"type" db:"type"`
	Status       NotificationStatus `json:"status" db:"status"`
	EmailContent string             `json:"email_content" db:"email_content"`
	SentAt       *time.Time         `json:"sent_at" db:"sent_at"`

	// Relationships
	Loan     Loan      `json:"loan,omitempty" belongs_to:"loans"`
	Investor *Investor `json:"investor,omitempty" belongs_to:"investors"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (n *Notification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(n.LoanID), Name: "LoanID"},
		&validators.StringIsPresent{Field: n.EmailContent, Name: "EmailContent"},
		&validators.StringInclusion{Field: string(n.Type), Name: "Type", List: []string{
			string(NotificationTypeAgreementLetter),
			string(NotificationTypeInvestmentConfirmation),
			string(NotificationTypeDisbursementNotice),
		}},
		&validators.StringInclusion{Field: string(n.Status), Name: "Status", List: []string{
			string(NotificationStatusPending),
			string(NotificationStatusSent),
			string(NotificationStatusFailed),
		}},
	), nil
}
