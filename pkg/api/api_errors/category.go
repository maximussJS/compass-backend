package api_errors

import "errors"

var (
	ErrorCategoryAlreadyExists = errors.New("category already exists")
	ErrorCategoryNotFound      = errors.New("category not found")
)
