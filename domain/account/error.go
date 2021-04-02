package account

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrAccountNotFound = errors.New("account not found")
	ErrUnauthorized = errors.New("unauthorized access: usernams or password incorrects")
)

