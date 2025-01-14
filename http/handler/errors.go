package handler

import (
	"errors"
)

var (
	ErrNilHandler = errors.New("handler cannot be nil")
	ErrNilRequest = errors.New("request cannot be nil")
)
