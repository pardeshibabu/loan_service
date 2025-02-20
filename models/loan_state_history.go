package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// LoanStateHistory model
type LoanStateHistory struct {
	Model
	LoanID      int64     `json:"loan_id" db:"loan_id"`
	ChangedByID int64     `json:"changed_by_id" db:"changed_by_id"`
	FromStatus  *string   `json:"from_status" db:"from_status"`
	ToStatus    string    `json:"to_status" db:"to_status"`
	ChangedAt   time.Time `json:"changed_at" db:"changed_at"`

	// Relationships
	Loan      Loan     `json:"loan,omitempty" belongs_to:"loans"`
	ChangedBy Employee `json:"changed_by,omitempty" belongs_to:"employees"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (l *LoanStateHistory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(l.LoanID), Name: "LoanID"},
		&validators.IntIsPresent{Field: int(l.ChangedByID), Name: "ChangedByID"},
		&validators.StringIsPresent{Field: l.ToStatus, Name: "ToStatus"},
		&validators.TimeIsPresent{Field: l.ChangedAt, Name: "ChangedAt"},
	), nil
}
