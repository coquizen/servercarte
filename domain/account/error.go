package account

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrNotAuthorized   = errors.New("unauthorized access: username or password incorrect")
)

