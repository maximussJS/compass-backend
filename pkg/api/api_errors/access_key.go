package api_errors

import "errors"

var (
	ErrorAccessKeyAlreadyExists = errors.New("access key already exists")
	ErrorAccessKeyNotFound      = errors.New("access key not found")
)
