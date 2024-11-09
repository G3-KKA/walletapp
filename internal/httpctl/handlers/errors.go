package handlers

import "errors"

var (
	ErrInvalidUUID            = errors.New("id is not valid UUID")
	ErrInernalError           = errors.New("internal error")
	ErrGotInvalidJSON         = errors.New("got invalid JSON")
	ErrUnsuccessfulValidation = errors.New("validation unsuccessful")
)
