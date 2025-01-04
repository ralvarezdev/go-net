package response

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
)

var (
	InternalServerError = NewJSONErrorResponseFromString(gonethttp.InternalServerError)
)
