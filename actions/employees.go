package actions

import (
	"fmt"
	"net/http"

	"loan_service/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// EmployeesList gets all employees
func EmployeesList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	employees := &[]models.Employee{}

	q := tx.PaginateFromParams(c.Params())

	// Add role filter if provided
	if role := c.Param("role"); role != "" {
		q = q.Where("role = ?", role)
	}

	// Add status filter if provided
	if status := c.Param("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.All(employees); err != nil {
		return err
	}

	c.Set("pagination", q.Paginator)
	return c.Render(http.StatusOK, r.JSON(employees))
}

// EmployeesShow gets the data for one employee
func EmployeesShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	employee := &models.Employee{}

	if err := tx.Find(employee, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(employee))
}

// EmployeesCreate adds an employee to the DB
func EmployeesCreate(c buffalo.Context) error {
	employee := &models.Employee{}

	if err := c.Bind(employee); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)

	// Set default status if not provided
	if employee.Status == "" {
		employee.Status = models.StatusActive
	}

	// Check if employee_id is unique
	exists, err := tx.Where("id = ?", employee.ID).Exists(&models.Employee{})
	if err != nil {
		return err
	}
	if exists {
		return c.Error(http.StatusBadRequest, fmt.Errorf("id must be unique"))
	}

	verrs, err := tx.ValidateAndCreate(employee)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(employee))
}

// EmployeesUpdate changes an employee in the DB
func EmployeesUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	employee := &models.Employee{}

	if err := tx.Find(employee, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	oldID := employee.ID
	if err := c.Bind(employee); err != nil {
		return err
	}

	// Check if employee_id is changed and is unique
	if oldID != employee.ID {
		exists, err := tx.Where("id = ?", employee.ID).Exists(&models.Employee{})
		if err != nil {
			return err
		}
		if exists {
			return c.Error(http.StatusBadRequest, fmt.Errorf("id must be unique"))
		}
	}

	verrs, err := tx.ValidateAndUpdate(employee)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(employee))
}

// EmployeesDelete deletes an employee from the DB
func EmployeesDelete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	employee := &models.Employee{}

	if err := tx.Find(employee, c.Param("id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check if employee is associated with any active loans
	count, err := tx.Where("(field_validator_id = ? OR field_officer_id = ?) AND status != ?",
		employee.ID, employee.ID, models.LoanStatusDisbursed).Count(&models.Loan{})
	if err != nil {
		return err
	}

	if count > 0 {
		return c.Error(http.StatusBadRequest, fmt.Errorf("cannot delete employee with active loans"))
	}

	if err := tx.Destroy(employee); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Employee deleted successfully"}))
}
