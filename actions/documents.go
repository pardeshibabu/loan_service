package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"loan_service/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// DocumentsList gets all documents for a loan
func DocumentsList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	documents := &[]models.Document{}

	// Get documents for specific loan
	q := tx.Where("loan_id = ?", c.Param("id"))
	q = tx.PaginateFromParams(c.Params())

	if err := q.All(documents); err != nil {
		return err
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.JSON(documents))
}

// DocumentsCreate creates a new document
func DocumentsCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	document := &models.Document{}

	// Bind document data
	if err := c.Bind(document); err != nil {
		return err
	}

	// Set loan ID from URL
	loanID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.Error(http.StatusBadRequest, fmt.Errorf("invalid loan ID"))
	}
	document.LoanID = loanID

	// Get the loan to verify it exists
	loan := &models.Loan{}
	if err := tx.Find(loan, loanID); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("loan not found"))
	}

	// Create document
	verrs, err := tx.ValidateAndCreate(document)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(document))
}
