package api_errors

import "errors"

var (
	ErrorUserAlreadyExists  = errors.New("user already exists")
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorInvalidToken       = errors.New("invalid token")
	ErrorUnauthorized       = errors.New("unauthorized")
)
