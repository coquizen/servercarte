package account

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrUnauthorized = errors.New("unauthorized access: username or password incorrect")
)

