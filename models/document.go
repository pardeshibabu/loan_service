package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeValidationProof DocumentType = "validation_proof"
	DocumentTypeAgreementLetter DocumentType = "agreement_letter"
	DocumentTypeSignedAgreement DocumentType = "signed_agreement"
	DocumentTypeKYC             DocumentType = "kyc"
)

// Document model
type Document struct {
	Model
	LoanID       int64        `json:"loan_id" db:"loan_id"`
	DocumentType DocumentType `json:"document_type" db:"document_type"`
	FileURL      string       `json:"file_url" db:"file_url"`
	UploadedByID int64        `json:"uploaded_by_id" db:"uploaded_by_id"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`

	// Relationships
	Loan       Loan     `json:"loan,omitempty" belongs_to:"loans"`
	UploadedBy Employee `json:"uploaded_by,omitempty" belongs_to:"employees"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(d.LoanID), Name: "LoanID"},
		&validators.StringIsPresent{Field: d.FileURL, Name: "FileURL"},
		&validators.IntIsPresent{Field: int(d.UploadedByID), Name: "UploadedByID"},
		&validators.StringInclusion{Field: string(d.DocumentType), Name: "DocumentType", List: []string{
			string(DocumentTypeValidationProof),
			string(DocumentTypeAgreementLetter),
			string(DocumentTypeSignedAgreement),
			string(DocumentTypeKYC),
		}},
	), nil
}
