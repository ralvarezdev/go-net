package json

import (
	"encoding/json"
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
)

type (
	// Decoder struct
	Decoder struct {
		mode             *goflagsmode.Flag
		responsesHandler gonethttpresponsehandler.ResponsesHandler
	}
)

// NewDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//   - responsesHandler: The HTTP handler
//
// Returns:
//
//   - *Decoder: The default decoder
//   - error: The error if any
func NewDecoder(
	mode *goflagsmode.Flag,
	responsesHandler gonethttpresponsehandler.ResponsesHandler,
) (*Decoder, error) {
	// Check if the handler is nil
	if responsesHandler == nil {
		return nil, gonethttpresponsehandler.ErrNilHandler
	}

	return &Decoder{
		mode,
		responsesHandler,
	}, nil
}

// Decode decodes the JSON request body and stores it in the destination
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	// Check the content type
	if !CheckContentType(r) {
		d.responsesHandler.HandleFieldFailResponseWithCode(
			w,
			ErrInvalidContentTypeField,
			ErrInvalidContentType,
			ErrCodeInvalidContentType,
			http.StatusUnsupportedMediaType,
		)
		return ErrInvalidContentType
	}

	// Check the decoder destination
	if dest == nil {
		d.responsesHandler.HandleDebugErrorResponseWithCode(
			w,
			ErrNilDestination,
			gonethttp.ErrInternalServerError,
			ErrCodeNilDestination,
			http.StatusInternalServerError,
		)
		return ErrNilDestination
	}

	// Get the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.responsesHandler.HandleDebugErrorResponseWithCode(
			w,
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeFailedToReadBody,
			http.StatusInternalServerError,
		)
		return err
	}

	// Decode JSON body into destination
	if err = json.Unmarshal(body, dest); err != nil {
		_ = BodyDecodeErrorHandler(w, err, d.responsesHandler)
	}
	return err
}
