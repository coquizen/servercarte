package user

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"

	"github.com/CaninoDev/gastro/server/domain"
)

var phoneRegExp = "(?:^|[^0-9])(1[34578][0-9]{9})(?:$|[^0-9])"

type User struct {
	domain.Base
	FirstName string `json:"first_name" gorm:"unique,not null"`
	LastName  string `json:"last_name,omitempty" gorm:"unique,null"`
	Address1  string `json:"address_1" gorm:"not null"`
	Address2  *string `json:"address_2,omitempty" gorm:"null"`
	ZipCode   uint   `json:"zip_code" gorm:"not null"`
	Email     string `json:"email" gorm:"unique,not null"`
	TelephoneNumber string `json:"phone,omitempty" gorm:"null"`
}

func (u *User) Validate() error {
	if u.FirstName == "" || u.LastName == "" {
		return errors.New("first name is empty")
	}
	if u.LastName == "" {
		return errors.New("last name is empty")
	}
	if hasNumber(u.FirstName) {
		return errors.New("first name has number character(s)")
	}
	if hasNumber(u.LastName) {
		return errors.New("last name has number character(s)")
	}
	if u.Address1 == "" {
		return errors.New("address1 is empty")
	}
	if !(len([]rune(strconv.Itoa(int(u.ZipCode)))) == 5 || len([]rune(strconv.Itoa(int(u.ZipCode)))) == 9)  {
		return errors.New("invalid zip code")
	}
	if len([]rune(u.TelephoneNumber)) > 0 {
		if ok, _ := regexp.Match(phoneRegExp, []byte(u.TelephoneNumber)); !ok {
		return errors.New("not a valid telephone number")
	}
	}

	return nil
}

func hasNumber(str string) bool {
	for _, character := range str {
		if unicode.IsNumber(character) {
			return true
		}
	}
	return false
}
