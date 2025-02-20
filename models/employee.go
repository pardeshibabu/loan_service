package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// EmployeeRole represents the role of an employee
type EmployeeRole string

const (
	RoleFieldValidator EmployeeRole = "field_validator"
	RoleFieldOfficer   EmployeeRole = "field_officer"
	RoleAdmin          EmployeeRole = "admin"
)

// EmployeeStatus represents the status of an employee
type EmployeeStatus string

const (
	StatusActive   EmployeeStatus = "active"
	StatusInactive EmployeeStatus = "inactive"
)

// Employee model
type Employee struct {
	Model
	EmployeeID string         `json:"employee_id" db:"employee_id"`
	FirstName  string         `json:"first_name" db:"first_name"`
	LastName   string         `json:"last_name" db:"last_name"`
	Email      string         `json:"email" db:"email"`
	Role       EmployeeRole   `json:"role" db:"role"`
	Status     EmployeeStatus `json:"status" db:"status"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (e *Employee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: e.EmployeeID, Name: "EmployeeID"},
		&validators.StringIsPresent{Field: e.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: e.LastName, Name: "LastName"},
		&validators.EmailIsPresent{Name: "Email", Field: e.Email},
		&validators.StringInclusion{Field: string(e.Role), Name: "Role", List: []string{
			string(RoleFieldValidator),
			string(RoleFieldOfficer),
			string(RoleAdmin),
		}},
		&validators.StringInclusion{Field: string(e.Status), Name: "Status", List: []string{
			string(StatusActive),
			string(StatusInactive),
		}},
	), nil
}
