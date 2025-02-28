package api_errors

import "errors"

var (
	ErrorUserPasswordIsTheSame = errors.New("password is the same")
	ErrorNameIsTheSame         = errors.New("name is the same")
	ErrorEmailAlreadyConfirmed = errors.New("email is already confirmed")
	ErrorUserNotFound          = errors.New("user not found")
)
