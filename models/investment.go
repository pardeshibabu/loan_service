package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// InvestmentStatus represents the status of an investment
type InvestmentStatus string

const (
	InvestmentStatusActive    InvestmentStatus = "active"
	InvestmentStatusCompleted InvestmentStatus = "completed"
	InvestmentStatusCancelled InvestmentStatus = "cancelled"
)

// Investment model
type Investment struct {
	Model
	LoanID             int64            `json:"loan_id" db:"loan_id"`
	InvestorID         int64            `json:"investor_id" db:"investor_id"`
	Amount             float64          `json:"amount" db:"amount"`
	AgreementLetterURL *string          `json:"agreement_letter_url" db:"agreement_letter_url"`
	Status             InvestmentStatus `json:"status" db:"status"`

	// Relationships
	Loan     Loan     `json:"loan,omitempty" belongs_to:"loans"`
	Investor Investor `json:"investor,omitempty" belongs_to:"investors"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (i *Investment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(i.LoanID), Name: "LoanID"},
		&validators.IntIsPresent{Field: int(i.InvestorID), Name: "InvestorID"},
		&Float64IsPresent{Field: i.Amount, Name: "Amount"},
		&validators.StringInclusion{Field: string(i.Status), Name: "Status", List: []string{
			string(InvestmentStatusActive),
			string(InvestmentStatusCompleted),
			string(InvestmentStatusCancelled),
		}},
	), nil
}
