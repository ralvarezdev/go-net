package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Response is the interface for the responses
	Response interface {
		Body(mode *goflagsmode.Flag) interface{}
		HTTPStatus() int
	}
)
