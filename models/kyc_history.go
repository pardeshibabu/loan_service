package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// KYCHistory tracks KYC status changes
type KYCHistory struct {
	Model
	InvestorID int64     `json:"investor_id" db:"investor_id"`
	FromStatus KYCStatus `json:"from_status" db:"from_status"`
	ToStatus   KYCStatus `json:"to_status" db:"to_status"`
	ReviewerID int64     `json:"reviewer_id" db:"reviewer_id"`
	Comments   string    `json:"comments" db:"comments"`
	ReviewedAt time.Time `json:"reviewed_at" db:"reviewed_at"`

	// Relationships
	Investor Investor `json:"investor,omitempty" belongs_to:"investors"`
	Reviewer Employee `json:"reviewer,omitempty" belongs_to:"employees"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (h *KYCHistory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(h.InvestorID), Name: "InvestorID"},
		&validators.IntIsPresent{Field: int(h.ReviewerID), Name: "ReviewerID"},
		&validators.StringIsPresent{Field: string(h.ToStatus), Name: "ToStatus"},
		&validators.TimeIsPresent{Field: h.ReviewedAt, Name: "ReviewedAt"},
	), nil
}
