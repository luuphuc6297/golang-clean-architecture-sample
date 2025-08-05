// Package validators provides validation functions for domain entities and request data.
// It includes validation for common fields like email, password, and business-specific rules.
package validators

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates that the provided email address is in a valid format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.ErrEmailIsRequired
	}
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidEmail
	}
	return nil
}

// ValidateRequired validates that a required field is not empty
func ValidateRequired(field, value string) error {
	if value == "" {
		switch field {
		case constants.FieldFirstName:
			return errors.ErrFirstNameIsRequired
		case constants.FieldLastName:
			return errors.ErrLastNameIsRequired
		case constants.FieldRole:
			return errors.ErrRoleIsRequired
		case constants.FieldName:
			return errors.ErrCategoryRequired
		default:
			return errors.ErrInvalidRequest
		}
	}
	return nil
}

// ValidatePrice validates that a price is positive and within reasonable limits
func ValidatePrice(price float64) error {
	if price <= 0 {
		return errors.ErrInvalidRequest
	}
	return nil
}

// ValidateStock validates that stock quantity is non-negative
func ValidateStock(stock int) error {
	if stock < 0 {
		return errors.ErrInvalidRequest
	}
	return nil
}

// ValidateRole validates that the role is one of the allowed values
func ValidateRole(role string) error {
	if role == "" {
		return errors.ErrRoleIsRequired
	}
	if role != constants.RoleUser && role != constants.RoleAdmin {
		return errors.ErrInvalidRole
	}
	return nil
}

// ValidatePassword validates that password meets minimum security requirements
func ValidatePassword(password string) error {
	if password == "" {
		return errors.ErrPasswordRequired
	}
	if len(password) < 6 {
		return errors.ErrPasswordTooShort
	}
	return nil
}

// ValidateRegisterRequest validates all fields required for user registration
func ValidateRegisterRequest(email, password, firstName, lastName string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if err := ValidateRequired(constants.FieldFirstName, firstName); err != nil {
		return err
	}
	if err := ValidateRequired(constants.FieldLastName, lastName); err != nil {
		return err
	}
	return nil
}

// ValidateLoginRequest validates fields required for user login
func ValidateLoginRequest(email, password string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if password == "" {
		return errors.ErrPasswordRequired
	}
	return nil
}