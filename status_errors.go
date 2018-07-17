package userengage

import "errors"

// Possible user engage errors
var (
	ErrUnauthorized = errors.New("unauthorized api key")
	ErrNotFound     = errors.New("resource not found")
	ErrServerError  = errors.New("internal server error")
)

var statusErrors = map[int]error{
	401: ErrUnauthorized,
	404: ErrNotFound,
	500: ErrServerError,
}
