package actions

import (
	"fmt"
	"net/http"
	"strconv"

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
	investment := &models.Investment{}

	// Bind investment data
	if err := c.Bind(investment); err != nil {
		return err
	}

	// Set loan ID from URL
	loanID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.Error(http.StatusBadRequest, fmt.Errorf("invalid loan ID"))
	}
	investment.LoanID = loanID

	// Get the loan to check status and amount
	loan := &models.Loan{}
	if err := tx.Find(loan, investment.LoanID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Validate loan status
	if loan.Status != models.LoanStatusApproved {
		return c.Error(http.StatusBadRequest, fmt.Errorf("loan must be in approved state for investment"))
	}

	// Calculate total invested amount
	var totalInvested float64
	if err := tx.Where("loan_id = ?", loan.ID).Select("SUM(amount)").First(&totalInvested); err != nil {
		return err
	}

	// Check if new investment would exceed principal
	if totalInvested+investment.Amount > loan.PrincipalAmount {
		return c.Error(http.StatusBadRequest, fmt.Errorf("investment would exceed loan principal amount"))
	}

	// Create investment
	investment.Status = models.InvestmentStatusActive
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

		// Create state history
		history := &models.LoanStateHistory{
			Loan:        *loan,
			ChangedByID: investment.InvestorID, // Using investor as the changer
			FromStatus:  strPtr(string(models.LoanStatusApproved)),
			ToStatus:    string(models.LoanStatusInvested),
			ChangedAt:   investment.CreatedAt,
		}

		if err := tx.Create(history); err != nil {
			return err
		}
	}

	return c.Render(http.StatusCreated, r.JSON(investment))
}
