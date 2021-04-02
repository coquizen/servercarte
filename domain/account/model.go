package account

import (
	"time"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/domain"
	"github.com/CaninoDev/gastro/server/domain/user"
)


type AccessLevel int

const (
	Admin AccessLevel = iota
	Employee
	Guest
)

// Account contains a user's account properties for interacting with the system
type Account struct {
	domain.Base
	UserID    uuid.UUID   `json:"user_id" gorm:"not null"`
	User      user.User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Username  string      `json:"username" gorm:"not null"`
	Password  string      `json:"password" gorm:"not null"`
	Role      AccessLevel `json:"role" gorm:"not null"`
	Token     string      `json:"token,omitempty" gorm:"null"`
	LastLogin time.Time   `json:"last_login,omitempty" gorm:"null"`
}


// NewAccountRequest represent the request struct for Create endpoint
type NewAccountRequest struct {
	FirstName       string      `json:"first_name"`
	LastName        string      `json:"last_name"`
	Address1        string      `json:"address_1"`
	Address2        *string     `json:"address_2,omitempty"`
	ZipCode         uint        `json:"zip_code"`
	Username        string      `json:"username"`
	Password        string      `json:"password"`
	Role            AccessLevel `json:"role"`
	PasswordConfirm string      `json:"password_confirm"`
	Email           string      `json:"email"`
}

// UpdateAccountRequest represent the request struct for Update endpoint
type UpdateAccountRequest struct {
	ID       uuid.UUID    `json:"-"`
	Address1 *string      `json:"address_1,omitempty"`
	Address2 *string      `json:"address_2,omitempty"`
	ZipCode  *uint        `json:"zip_code,omitempty"`
	Email    *string      `json:"email,omitempty"`
	Role     *AccessLevel `json:"role,omitempty"`
}



