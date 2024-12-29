package response

import (
	"errors"
)

var (
	ErrNilHandler = errors.New("response handler cannot be nil")
)
