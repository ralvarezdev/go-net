package response

import (
	gonethttperrors "github.com/ralvarezdev/go-net/http/errors"
	"net/http"
)

var (
	InternalServerError = NewErrorResponse(
		gonethttperrors.InternalServerError,
		nil,
		nil,
		http.StatusInternalServerError,
	)
)
