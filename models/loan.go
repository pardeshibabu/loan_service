package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// LoanStatus represents the status of a loan
type LoanStatus string

const (
	LoanStatusProposed  LoanStatus = "proposed"
	LoanStatusApproved  LoanStatus = "approved"
	LoanStatusInvested  LoanStatus = "invested"
	LoanStatusDisbursed LoanStatus = "disbursed"
)

// Loan model
type Loan struct {
	ID                 int64      `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	BorrowerID         int64      `json:"borrower_id" db:"borrower_id"`
	PrincipalAmount    float64    `json:"principal_amount" db:"principal_amount"`
	Rate               float64    `json:"rate" db:"rate"`
	ROI                float64    `json:"roi" db:"roi"`
	Status             LoanStatus `json:"status" db:"status"`
	ApprovalDate       *time.Time `json:"approval_date" db:"approval_date"`
	FieldValidatorID   *int64     `json:"field_validator_id" db:"field_validator_id"`
	ValidationProofURL *string    `json:"validation_proof_url" db:"validation_proof_url"`
	DisbursementDate   *time.Time `json:"disbursement_date" db:"disbursement_date"`
	FieldOfficerID     *int64     `json:"field_officer_id" db:"field_officer_id"`
	SignedAgreementURL *string    `json:"signed_agreement_url" db:"signed_agreement_url"`
	AgreementLetterURL *string    `json:"agreement_letter_url" db:"agreement_letter_url"`

	// Relationships
	Borrower       *Borrower    `json:"borrower" belongs_to:"borrowers" fk_id:"borrower_id"`
	FieldValidator *Employee    `json:"field_validator" belongs_to:"employees" fk_id:"field_validator_id"`
	FieldOfficer   *Employee    `json:"field_officer" belongs_to:"employees" fk_id:"field_officer_id"`
	Investments    []Investment `json:"investments,omitempty" has_many:"investments"`
	Documents      []Document   `json:"documents,omitempty" has_many:"documents"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (l *Loan) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(l.BorrowerID), Name: "BorrowerID"},
		&Float64IsPresent{Field: l.PrincipalAmount, Name: "PrincipalAmount"},
		&Float64IsPresent{Field: l.Rate, Name: "Rate"},
		&Float64IsPresent{Field: l.ROI, Name: "ROI"},
		&validators.StringInclusion{Field: string(l.Status), Name: "Status", List: []string{
			string(LoanStatusProposed),
			string(LoanStatusApproved),
			string(LoanStatusInvested),
			string(LoanStatusDisbursed),
		}},
	), nil
}
