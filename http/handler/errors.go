package handler

import (
	"errors"
)

var (
	ErrNilHandler = errors.New("response handler cannot be nil")
)
