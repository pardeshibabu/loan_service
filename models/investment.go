package models

import (
	"time"

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
	ID             int64            `json:"id" db:"id"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
	LoanID         int64            `json:"loan_id" db:"loan_id"`
	InvestorID     int64            `json:"investor_id" db:"investor_id"`
	Amount         float64          `json:"amount" db:"amount"`
	Status         InvestmentStatus `json:"status" db:"status"`
	InvestmentDate *time.Time       `json:"investment_date" db:"investment_date"`

	// Relationships
	Loan     *Loan     `json:"loan" belongs_to:"loans" fk_id:"loan_id"`
	Investor *Investor `json:"investor" belongs_to:"investors" fk_id:"investor_id"`
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
