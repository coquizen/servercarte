package model

import (
	"time"

	"github.com/google/uuid"
)


type AccessLevel int

const (
	Admin AccessLevel = iota
	Employee
	Guest
)

// Account contains a user's account properties for interacting with the system
type Account struct {
	Base
	UserID      uuid.UUID   `json:"userID" gorm:"not null"`
	User        User        `comment:"account BELONGS TO user"`
	Username    string      `json:"username" gorm:"not null"`
	Password    string      `json:"password" gorm:"not null"`
	Role        AccessLevel `json:"role" gorm:"not null"`
	Token       string      `json:"token,omitempty" gorm:"null"`
	TokenExpiry time.Time   `json:"token_expiry,omitempty" gorm:"null"`
	LastLogin   time.Time   `json:"last_login,omitempty" gorm:"null"`
}

