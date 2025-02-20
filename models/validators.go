package models

import (
	"errors"
	"regexp"

	"github.com/gobuffalo/validate/v3"
)

// Float64IsPresent is a validator that checks if a float64 is not zero
type Float64IsPresent struct {
	Name  string
	Field float64
}

// IsValid adds an error if the field is zero
func (v *Float64IsPresent) IsValid(errors *validate.Errors) {
	if v.Field == 0 {
		errors.Add(v.Name, "Must be present")
	}
}

// Email validates an email string
func Email(email string) error {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Must be a valid email address")
	}
	return nil
}
