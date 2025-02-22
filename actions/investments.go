package actions

import (
	"fmt"
	"net/http"

	"loan_service/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// InvestmentsList gets all investments for a loan
func InvestmentsList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investments := &[]models.Investment{}

	// Get investments for specific loan
	q := tx.Where("loan_id = ?", c.Param("id"))
	q = tx.PaginateFromParams(c.Params())

	if err := q.All(investments); err != nil {
		return err
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.JSON(investments))
}

// InvestmentsShow gets a specific investment
func InvestmentsShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investment := &models.Investment{}

	if err := tx.Find(investment, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(investment))
}

// InvestmentsCreate creates a new investment
func InvestmentsCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Get the loan first
	loan := &models.Loan{}
	if err := tx.Find(loan, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("loan not found"))
	}

	// Validate loan status
	if loan.Status != models.LoanStatusApproved {
		return c.Error(http.StatusBadRequest, fmt.Errorf("loan must be in approved state for investment"))
	}

	// Create investment
	investment := &models.Investment{
		LoanID: loan.ID,
		Status: models.InvestmentStatusActive,
	}

	// Bind request data
	if err := c.Bind(investment); err != nil {
		return err
	}

	// Calculate total invested amount
	var totalInvested float64
	query := tx.RawQuery("SELECT COALESCE(SUM(amount), 0) FROM investments WHERE loan_id = ?", loan.ID)
	if err := query.First(&totalInvested); err != nil {
		return err
	}

	// Calculate remaining amount
	remainingAmount := loan.PrincipalAmount - totalInvested

	// Check if new investment would exceed principal
	if investment.Amount > remainingAmount {
		return c.Error(http.StatusBadRequest, fmt.Errorf(
			"investment amount %.2f exceeds available amount %.2f (total loan amount: %.2f, already invested: %.2f)",
			investment.Amount, remainingAmount, loan.PrincipalAmount, totalInvested,
		))
	}

	// Create investment
	verrs, err := tx.ValidateAndCreate(investment)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	// Check if loan is fully invested
	if totalInvested+investment.Amount == loan.PrincipalAmount {
		loan.Status = models.LoanStatusInvested
		if err := tx.Update(loan); err != nil {
			return err
		}
	}

	// Fetch fresh investment with relationships
	if err := tx.Eager("Loan", "Investor").Find(investment, investment.ID); err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(investment))
}
