package authentication

import "errors"
var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrExpiredToken = errors.New("token expired")
)
