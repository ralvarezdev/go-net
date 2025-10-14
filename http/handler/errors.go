package handler

import (
	"errors"
)

var (
	ErrNilHandler = errors.New("requests and responses handler cannot be nil")
)
