package model

// User struct represent the user that will have admin access.
type User struct {
	Base
	FirstName string `json:"first_name" gorm:"unique,not null"`
	LastName  string `json:"last_name,omitempty" gorm:"unique,null"`
	Address1  string `json:"address_1" gorm:"not null"`
	Address2  string `json:"address_2,omitempty" gorm:"null"`
	ZipCode   uint   `json:"zip_code" gorm:"not null"`
	Email     string `json:"email" gorm:"unique,not null"`
}
