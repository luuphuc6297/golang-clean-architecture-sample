package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

type TestHelper struct{}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (th *TestHelper) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (th *TestHelper) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
