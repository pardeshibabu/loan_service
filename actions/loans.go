package actions

import (
	"fmt"
	"net/http"
	"time"

	"loan_service/models"
	"loan_service/services/notification"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// LoansList default implementation.
func LoansList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	loans := &[]models.Loan{}

	// Paginate results. Params "page" and "per_page" control pagination.
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Loans from the DB
	if err := q.All(loans); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(loans))
}

// LoansShow default implementation.
func LoansShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	loan := &models.Loan{}

	if err := tx.Find(loan, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(loan))
}

// LoansCreate creates a new loan
func LoansCreate(c buffalo.Context) error {
	loan := &models.Loan{}
	if err := c.Bind(loan); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	loan.Status = models.LoanStatusProposed

	verrs, err := tx.ValidateAndCreate(loan)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	// Fetch the fresh loan with relationships
	if err := tx.Eager("Borrower").Find(loan, loan.ID); err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(loan))
}

// LoansApprove handles loan approval
func LoansApprove(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	loan := &models.Loan{}

	// Find the loan with eager loading
	if err := tx.Eager("Borrower", "FieldValidator", "FieldOfficer").Find(loan, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Validate current status
	if loan.Status != models.LoanStatusProposed {
		return c.Error(http.StatusBadRequest, fmt.Errorf("loan must be in proposed state to approve"))
	}

	// Bind approval data
	type ApprovalRequest struct {
		FieldValidatorID   int64  `json:"field_validator_id"`
		ValidationProofURL string `json:"validation_proof_url"`
	}

	req := &ApprovalRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	// Update loan
	now := time.Now()
	loan.Status = models.LoanStatusApproved
	loan.ApprovalDate = &now
	loan.FieldValidatorID = &req.FieldValidatorID
	loan.ValidationProofURL = &req.ValidationProofURL

	verrs, err := tx.ValidateAndUpdate(loan)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	// Create state history
	history := &models.LoanStateHistory{
		LoanID:      loan.ID,
		ChangedByID: req.FieldValidatorID,
		FromStatus:  strPtr(string(models.LoanStatusProposed)),
		ToStatus:    string(models.LoanStatusApproved),
		ChangedAt:   now,
	}

	if err := tx.Create(history); err != nil {
		return err
	}

	// After updating, fetch the fresh loan with relationships
	if err := tx.Eager("Borrower", "FieldValidator", "FieldOfficer").Find(loan, loan.ID); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(loan))
}

// LoansDisburse handles loan disbursement
func LoansDisburse(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	loan := &models.Loan{}

	// Find the loan
	if err := tx.Find(loan, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Validate current status
	if loan.Status != models.LoanStatusInvested {
		return c.Error(http.StatusBadRequest, fmt.Errorf("loan must be in invested state to disburse"))
	}

	// Bind disbursement data
	type DisburseRequest struct {
		FieldOfficerID     int64  `json:"field_officer_id"`
		SignedAgreementURL string `json:"signed_agreement_url"`
	}

	req := &DisburseRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	// Update loan
	now := time.Now()
	loan.Status = models.LoanStatusDisbursed
	loan.DisbursementDate = &now
	loan.FieldOfficerID = &req.FieldOfficerID
	loan.SignedAgreementURL = &req.SignedAgreementURL

	verrs, err := tx.ValidateAndUpdate(loan)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	// Create state history
	history := &models.LoanStateHistory{
		LoanID:      loan.ID,
		ChangedByID: req.FieldOfficerID,
		FromStatus:  strPtr(string(models.LoanStatusInvested)),
		ToStatus:    string(models.LoanStatusDisbursed),
		ChangedAt:   now,
	}

	if err := tx.Create(history); err != nil {
		return err
	}

	// Send disbursement notifications
	notificationService, err := notification.NewEmailService(tx)
	if err != nil {
		// Log error but continue
		c.Logger().Error(err)
	} else {
		if err := notificationService.SendDisbursementNotice(c, loan); err != nil {
			// Log error but don't fail the request
			c.Logger().Error(err)
		}
	}

	return c.Render(http.StatusOK, r.JSON(loan))
}
