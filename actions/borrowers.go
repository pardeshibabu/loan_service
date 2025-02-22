package actions

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("Step 1: Starting borrower creation")

	// Initialize borrower struct
	borrower := &models.Borrower{}
	fmt.Println("Step 2: Initialized empty borrower struct")

	// Bind request data
	if err := c.Bind(borrower); err != nil {
		fmt.Printf("Error in binding data: %v\n", err)
		return err
	}
	debugData, _ := json.Marshal(borrower)
	fmt.Printf("Step 3: Bound request data: %s\n", string(debugData))

	// Get DB transaction
	fmt.Println("Step 4: Getting DB transaction")
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		fmt.Println("Error: Could not get DB connection from context")
		return fmt.Errorf("database connection not found")
	}
	fmt.Println("Step 5: Got DB transaction successfully")

	// Validate and create
	fmt.Println("Step 6: Starting ValidateAndCreate")
	fmt.Printf("DB Transaction: %+v\n", tx)
	fmt.Printf("Borrower Data: %+v\n", borrower)

	verrs, err := tx.ValidateAndCreate(borrower)
	if err != nil {
		fmt.Printf("Error in ValidateAndCreate: %v\n", err)
		// Print the full error details
		fmt.Printf("Full error: %+v\n", err)
		return err
	}
	fmt.Println("Step 7: ValidateAndCreate completed")

	// Check validation errors
	if verrs.HasAny() {
		fmt.Printf("Validation errors found: %v\n", verrs.String())
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	fmt.Println("Step 8: Borrower created successfully")
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
