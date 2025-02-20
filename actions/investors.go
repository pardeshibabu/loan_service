package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"loan_service/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// InvestorsList godoc
// @Summary List all investors
// @Description Get a paginated list of investors with optional filters
// @Tags investors
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param kyc_status query string false "KYC status filter" Enums(pending,approved,rejected)
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {array} models.Investor
// @Router /investors [get]
func InvestorsList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investors := &[]models.Investor{}

	// Parse search params
	search := SearchParams{}
	if err := c.Bind(&search); err != nil {
		return err
	}

	// Build query
	q := tx.Q()
	q = ApplySearch(q, search)

	// Add KYC status filter
	if status := c.Param("kyc_status"); status != "" {
		q = q.Where("kyc_status = ?", status)
	}

	// Add investment amount filter
	if minAmount := c.Param("min_investment"); minAmount != "" {
		q = q.Where("total_investment_amount >= ?", minAmount)
	}

	// Apply pagination
	q = q.PaginateFromParams(c.Params())

	if err := q.All(investors); err != nil {
		return err
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.JSON(investors))
}

// InvestorsShow gets the data for one investor
func InvestorsShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investor := &models.Investor{}

	if err := tx.Find(investor, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(investor))
}

// InvestorsCreate adds an investor to the DB
func InvestorsCreate(c buffalo.Context) error {
	investor := &models.Investor{}

	if err := c.Bind(investor); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)

	// Set default KYC status
	if investor.KYCStatus == "" {
		investor.KYCStatus = models.KYCStatusPending
	}

	verrs, err := tx.ValidateAndCreate(investor)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(investor))
}

// InvestorsUpdate changes an investor in the DB
func InvestorsUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investor := &models.Investor{}

	if err := tx.Find(investor, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(investor); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(investor)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(investor))
}

// InvestorsDelete deletes an investor from the DB
func InvestorsDelete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investor := &models.Investor{}

	if err := tx.Find(investor, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check if investor has any active investments
	count, err := tx.Where("investor_id = ? AND status = ?", investor.ID, models.InvestmentStatusActive).Count(&models.Investment{})
	if err != nil {
		return err
	}

	if count > 0 {
		return c.Error(http.StatusBadRequest, fmt.Errorf("cannot delete investor with active investments"))
	}

	if err := tx.Destroy(investor); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Investor deleted successfully"}))
}

// KYCUpdateRequest represents the request body for KYC update
type KYCUpdateRequest struct {
	Status     models.KYCStatus `json:"status"`
	Documents  []string         `json:"documents"`
	ReviewerID int64            `json:"reviewer_id"`
	Comments   string           `json:"comments"`
}

// InvestorsKYCUpdate godoc
// @Summary Update investor KYC status
// @Description Update an investor's KYC status and documents
// @Tags investors
// @Accept json
// @Produce json
// @Param id path int true "Investor ID"
// @Param kyc body KYCUpdateRequest true "KYC update info"
// @Success 200 {object} models.Investor
// @Router /investors/{id}/kyc [put]
func InvestorsKYCUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	investor := &models.Investor{}

	if err := tx.Find(investor, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Parse request
	req := &KYCUpdateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	// Validate status transition
	if !isValidKYCTransition(investor.KYCStatus, req.Status) {
		return c.Error(http.StatusBadRequest, fmt.Errorf("invalid KYC status transition"))
	}

	// Update investor
	investor.KYCStatus = req.Status
	docs := strings.Join(req.Documents, ",")
	investor.KYCDocuments = &docs

	verrs, err := tx.ValidateAndUpdate(investor)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	// Create KYC history record
	history := &models.KYCHistory{
		InvestorID: investor.ID,
		FromStatus: investor.KYCStatus,
		ToStatus:   req.Status,
		ReviewerID: req.ReviewerID,
		Comments:   req.Comments,
		ReviewedAt: time.Now(),
	}

	if err := tx.Create(history); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(investor))
}

// InvestmentSummary represents an investor's investment summary
type InvestmentSummary struct {
	TotalInvested       float64             `json:"total_invested"`
	ActiveInvestments   int                 `json:"active_investments"`
	TotalReturns        float64             `json:"total_returns"`
	InvestmentsByStatus map[string]float64  `json:"investments_by_status"`
	RecentInvestments   []models.Investment `json:"recent_investments"`
}

// InvestorsInvestmentSummary gets investment summary for an investor
func InvestorsInvestmentSummary(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Get all investments for the investor
	investments := []models.Investment{}
	if err := tx.Where("investor_id = ?", c.Param("id")).
		Order("created_at desc").All(&investments); err != nil {
		return err
	}

	// Calculate summary
	summary := InvestmentSummary{
		InvestmentsByStatus: make(map[string]float64),
	}

	for _, inv := range investments {
		summary.TotalInvested += inv.Amount

		if inv.Status == models.InvestmentStatusActive {
			summary.ActiveInvestments++
		}

		summary.InvestmentsByStatus[string(inv.Status)] += inv.Amount
	}

	// Get recent investments (last 5)
	if len(investments) > 5 {
		summary.RecentInvestments = investments[:5]
	} else {
		summary.RecentInvestments = investments
	}

	// Calculate returns
	var returns float64
	err := tx.RawQuery(`
		SELECT COALESCE(SUM(i.amount * l.roi / 100), 0) as total_returns 
		FROM investments i 
		JOIN loans l ON i.loan_id = l.id 
		WHERE i.investor_id = ? AND i.status = ?`,
		c.Param("id"), models.InvestmentStatusCompleted).First(&returns)

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	summary.TotalReturns = returns

	return c.Render(http.StatusOK, r.JSON(summary))
}

// Helper function to validate KYC status transitions
func isValidKYCTransition(from, to models.KYCStatus) bool {
	validTransitions := map[models.KYCStatus][]models.KYCStatus{
		models.KYCStatusPending: {
			models.KYCStatusApproved,
			models.KYCStatusRejected,
		},
		models.KYCStatusRejected: {
			models.KYCStatusPending,
		},
		models.KYCStatusApproved: {
			models.KYCStatusRejected,
		},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
