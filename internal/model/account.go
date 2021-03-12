package model

import "github.com/google/uuid"

type Account struct {
	Base
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	UserID uuid.UUID `json:"userID" gorm:"not null"`
	Role string `json:"role" gorm:"not null"`
}