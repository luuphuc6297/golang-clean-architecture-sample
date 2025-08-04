package entities

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/validators"
)

// User represents a user entity in the system
type User struct {
	BaseEntity
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Role      string `json:"role" gorm:"default:user"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
}

// TableName returns the database table name for the User entity
func (User) TableName() string {
	return "users"
}

// Validate validates the user entity fields and returns an error if invalid
func (u *User) Validate() error {
	if u.Role == "" {
		u.Role = constants.RoleUser
	}

	if err := validators.ValidateRole(u.Role); err != nil {
		return err
	}

	return nil
}

// IsAdmin returns true if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == constants.RoleAdmin
}
