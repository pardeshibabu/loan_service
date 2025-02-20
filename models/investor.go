package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// KYCStatus represents the KYC status of an investor
type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusApproved KYCStatus = "approved"
	KYCStatusRejected KYCStatus = "rejected"
)

// Investor model
type Investor struct {
	Model
	FirstName             string    `json:"first_name" db:"first_name"`
	LastName              string    `json:"last_name" db:"last_name"`
	Email                 string    `json:"email" db:"email"`
	Phone                 *string   `json:"phone" db:"phone"`
	KYCStatus             KYCStatus `json:"kyc_status" db:"kyc_status"`
	KYCDocuments          *string   `json:"kyc_documents" db:"kyc_documents"`
	TotalInvestmentAmount float64   `json:"total_investment_amount" db:"total_investment_amount"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (i *Investor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: i.LastName, Name: "LastName"},
		&validators.EmailIsPresent{Name: "Email", Field: i.Email},
		&validators.StringInclusion{Field: string(i.KYCStatus), Name: "KYCStatus", List: []string{
			string(KYCStatusPending),
			string(KYCStatusApproved),
			string(KYCStatusRejected),
		}},
	), nil
}
