package api_errors

import "errors"

var (
	ErrorFilePathNotFound = errors.New("file path not found")
	ErrorFileNotFound     = errors.New("file not found")
	ErrorFilesNotFound    = errors.New("files not found")
)
