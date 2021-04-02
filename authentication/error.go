package authentication

import "errors"
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrExpiredToken = errors.New("token expired")
	ErrNoValidTokenFound = errors.New("no valid token found in header")
)
