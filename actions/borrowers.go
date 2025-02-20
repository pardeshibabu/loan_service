package actions

import (
	"net/http"

	"loan_service/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// BorrowersList gets all borrowers
func BorrowersList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	borrowers := &[]models.Borrower{}

	// Parse search params
	search := SearchParams{}
	if err := c.Bind(&search); err != nil {
		return err
	}

	// Build query
	q := tx.Q()
	q = ApplySearch(q, search)

	// Apply pagination
	q = q.PaginateFromParams(c.Params())

	if err := q.All(borrowers); err != nil {
		return err
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.JSON(borrowers))
}

// BorrowersShow gets the data for one borrower
func BorrowersShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	borrower := &models.Borrower{}

	if err := tx.Find(borrower, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(borrower))
}

// BorrowersCreate adds a borrower to the DB
func BorrowersCreate(c buffalo.Context) error {
	borrower := &models.Borrower{}

	if err := c.Bind(borrower); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(borrower)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(borrower))
}

// BorrowersUpdate changes a borrower in the DB
func BorrowersUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	borrower := &models.Borrower{}

	if err := tx.Find(borrower, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(borrower); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(borrower)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(borrower))
}

// BorrowersDelete deletes a borrower from the DB
func BorrowersDelete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	borrower := &models.Borrower{}

	if err := tx.Find(borrower, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(borrower); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Borrower deleted successfully"}))
}
