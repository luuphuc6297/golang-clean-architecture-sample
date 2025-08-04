package validators

import (
	"regexp"

	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.ErrEmailIsRequired
	}
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidEmail
	}
	return nil
}

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

func ValidatePrice(price float64) error {
	if price <= 0 {
		return errors.ErrInvalidRequest
	}
	return nil
}

func ValidateStock(stock int) error {
	if stock < 0 {
		return errors.ErrInvalidRequest
	}
	return nil
}

func ValidateRole(role string) error {
	if role == "" {
		return errors.ErrRoleIsRequired
	}
	if role != constants.RoleUser && role != constants.RoleAdmin {
		return errors.ErrInvalidRole
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.ErrPasswordRequired
	}
	if len(password) < 6 {
		return errors.ErrPasswordTooShort
	}
	return nil
}

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

func ValidateLoginRequest(email, password string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if password == "" {
		return errors.ErrPasswordRequired
	}
	return nil
}
