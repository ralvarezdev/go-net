package response

import (
	"github.com/ralvarezdev/go-net/http/response"
	gonethttperrors "github.com/ralvarezdev/go-net/http/status/errors"
	"net/http"
)

var (
	InternalServerError = response.NewErrorResponse(
		gonethttperrors.InternalServerError,
		nil,
		nil,
		http.StatusInternalServerError,
	)
)
