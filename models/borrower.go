package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Borrower model
type Borrower struct {
	Model
	IDNumber  string  `json:"id_number" db:"id_number"`
	FirstName string  `json:"first_name" db:"first_name"`
	LastName  string  `json:"last_name" db:"last_name"`
	Email     *string `json:"email" db:"email"`
	Phone     *string `json:"phone" db:"phone"`
	Address   *string `json:"address" db:"address"`
	// We'll uncomment Loans after creating the Loan model
	// Loans     []Loan  `json:"loans,omitempty" has_many:"loans"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (b *Borrower) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var verrs validate.Errors

	if b.Email != nil {
		if err := Email(*b.Email); err != nil {
			verrs.Add("Email", err.Error())
		}
	}

	validationErrs := validate.Validate(
		&validators.StringIsPresent{Field: b.IDNumber, Name: "IDNumber"},
		&validators.StringIsPresent{Field: b.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: b.LastName, Name: "LastName"},
	)

	verrs.Append(validationErrs)
	return &verrs, nil
}
