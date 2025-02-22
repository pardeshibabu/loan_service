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

// mysql> desc employees;
// +------------+-------------------------------------------------+------+-----+---------+----------------+
// | Field      | Type                                            | Null | Key | Default | Extra          |
// +------------+-------------------------------------------------+------+-----+---------+----------------+
// | id         | bigint                                          | NO   | PRI | NULL    | auto_increment |
// | first_name | varchar(100)                                    | NO   |     | NULL    |                |
// | last_name  | varchar(100)                                    | NO   |     | NULL    |                |
// | email      | varchar(255)                                    | NO   |     | NULL    |                |
// | role       | enum('field_validator','field_officer','admin') | NO   |     | NULL    |                |
// | status     | enum('active','inactive')                       | YES  |     | active  |                |
// | created_at | datetime                                        | NO   |     | NULL    |                |
// | updated_at | datetime                                        | NO   |     | NULL    |                |
// +------------+-------------------------------------------------+------+-----+---------+----------------+
// 8 rows in set (0.00 sec)

// mysql>
// Employee model
type Employee struct {
	Model
	ID        int64          `json:"id" db:"id"`
	FirstName string         `json:"first_name" db:"first_name"`
	LastName  string         `json:"last_name" db:"last_name"`
	Email     string         `json:"email" db:"email"`
	Role      EmployeeRole   `json:"role" db:"role"`
	Status    EmployeeStatus `json:"status" db:"status"`
}

// Validate gets run every time you call a "pop.Validate*" method.
func (e *Employee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
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
