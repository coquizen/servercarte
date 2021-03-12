package model

// User struct represent the user that will have admin access.
type User struct {
	Base
	FirstName string  `json:"first_name,omitempty" gorm:"unique,not null"`
	LastName  string  `json:"last_name,omitempty" gorm:"unique,not null"`
	Addr   string  `json:"address" gorm:"not null"`
	ZipCode   uint    `json:"zip_code" gorm:"not null"`
	Email     string  `json:"email,omitempty" gorm:"unique,not null"`
	Account   Account `json:"-" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// // BeforeSave is a gorm hook that runs before a new row is inserted into the
// // db; In particular, hash the password before inserting into the
// // db.
// func (u *User) BeforeSave(tx *gorm.DB) error {
// 	passwordHash, err := security.Hash(u.Password)
// 	if err != nil {
// 		return err
// 	}
// 	u.Password = string(passwordHash)
// 	return nil
// }
