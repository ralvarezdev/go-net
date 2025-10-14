package json

import (
	"encoding/json"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
)

type (
	// StreamDecoder is the JSON decoder struct
	StreamDecoder struct {
		mode             *goflagsmode.Flag
		responsesHandler gonethttpresponsehandler.ResponsesHandler
	}
)

// NewStreamDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//   - responsesHandler: The HTTP handler to handle errors
//
// Returns:
//
//   - *StreamDecoder: The default decoder
//   - error: The error if any
func NewStreamDecoder(
	mode *goflagsmode.Flag,
	responsesHandler gonethttpresponsehandler.ResponsesHandler,
) (*StreamDecoder, error) {
	// Check if the response handler is nil
	if responsesHandler == nil {
		return nil, gonethttpresponsehandler.ErrNilHandler
	}

	return &StreamDecoder{
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
func (s StreamDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	// Check the content type
	if !CheckContentType(r) {
		s.responsesHandler.HandleFieldFailResponseWithCode(
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
		s.responsesHandler.HandleDebugErrorResponseWithCode(
			w,
			ErrNilDestination,
			gonethttp.ErrInternalServerError,
			ErrCodeNilDestination,
			http.StatusInternalServerError,
		)
		return ErrNilDestination
	}

	// Create a new reader from the body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode JSON body into destination
	if err = decoder.Decode(dest); err != nil {
		_ = BodyDecodeErrorHandler(w, err, s.responsesHandler)
	}
	return err
}
