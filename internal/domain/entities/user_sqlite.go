package entities

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/validators"

	"github.com/google/uuid"
)

// UserSQLite represents a user entity for SQLite database
type UserSQLite struct {
	BaseSQLiteEntity
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Role      string `json:"role" gorm:"default:user"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
}

// TableName returns the table name for UserSQLite entity
func (UserSQLite) TableName() string {
	return "users"
}

// Validate validates the SQLite user entity fields
func (u *UserSQLite) Validate() error {
	if err := validators.ValidateEmail(u.Email); err != nil {
		return err
	}
	if err := validators.ValidateRequired(constants.FieldFirstName, u.FirstName); err != nil {
		return err
	}
	if err := validators.ValidateRequired(constants.FieldLastName, u.LastName); err != nil {
		return err
	}
	if u.Role == "" {
		u.Role = constants.RoleUser
	}
	if err := validators.ValidateRole(u.Role); err != nil {
		return err
	}
	return nil
}

// IsAdmin returns true if the SQLite user has admin role
func (u *UserSQLite) IsAdmin() bool {
	return u.Role == constants.RoleAdmin
}

// ToUser converts SQLite user to domain user
func (u *UserSQLite) ToUser() *User {
	id, _ := uuid.Parse(u.ID)
	user := &User{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: u.DeletedAt,
		},
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
	}
	return user
}

// FromUser converts domain user to SQLite user
func FromUser(user *User) *UserSQLite {
	return &UserSQLite{
		BaseSQLiteEntity: BaseSQLiteEntity{
			ID:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt,
		},
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}
}
