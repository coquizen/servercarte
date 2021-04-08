package helpers

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// HasMixedCase verifies if a string has both upper and lower case characters.
func HasMixedCase(str string) bool {
	var hasLower, hasUpper bool
	hasLower = false
	hasUpper = false
	for _, character := range str {
		if unicode.IsLetter(character) {
			if !hasLower {
				hasLower = unicode.IsLower(character)
			}
			if !hasUpper {
				hasUpper = unicode.IsUpper(character)
			}
		}
		if hasLower && hasUpper {
			return true
		}
	}
	return false
}

// HasSpecialChar checks to see if a string has any special (printable) characters.
func HasSpecialChar(str string) bool {
	var specialChar = "!@#$%^&*"
	for _, character := range str {
		for _, spcCharacter := range specialChar {
			if character == spcCharacter {
				return true
			}
		}
	}
	return false
}

// HasAlphaNum checks to see if a string has both letters and numbers.
func HasAlphaNum(str string) bool {
	var hasAlpha, hasNumber bool
	hasAlpha = false
	hasNumber = false

	for _, character := range str {
		if unicode.IsLetter(character) {
			hasAlpha = true
		}
		if unicode.IsNumber(character) {
			hasNumber = true
		}
		if hasAlpha && hasNumber {
			return true
		}
	}
	return false
}

func IsEmailFormat(str string) bool {
	if len(str) < 3 && len(str) > 254 {
		return false
	}
	return emailRegex.MatchString(str)
}
