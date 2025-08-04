package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

// TestHelper provides utility functions for testing
type TestHelper struct{}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// HashPassword creates a bcrypt hash for testing purposes
func (th *TestHelper) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a password with its hash
func (th *TestHelper) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
