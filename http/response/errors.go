package response

import (
	"errors"
)

var (
	ErrNilResponse     = errors.New("response cannot be nil")
	ErrNilResponseCode = errors.New("response code cannot be nil")
)
