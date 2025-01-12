package validator

import (
	"errors"
)

var (
	ErrNilFailHandler = errors.New("fail handler cannot be nil")
)
