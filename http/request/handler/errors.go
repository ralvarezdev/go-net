package handler

import (
	"errors"
)

var (
	ErrCodeValidationFailed string
)

var (
	ErrNilHandler = errors.New("request handler cannot be nil")
)
