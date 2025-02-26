package api_errors

import "errors"

var (
	ErrorInvalidRequestParameter = errors.New("invalid request parameter")
	ErrorInternalServerError     = errors.New("internal server error")
)
