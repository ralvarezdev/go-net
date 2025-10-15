package json

import (
	"encoding/json"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// StreamDecoder is the JSON decoder struct
	StreamDecoder struct {
		mode *goflagsmode.Flag
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
func NewStreamDecoder(
	mode *goflagsmode.Flag,
) *StreamDecoder {
	return &StreamDecoder{
		mode,
	}
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
) error {
	// Check the content type
	if !CheckContentType(r) {
		return gonethttpresponse.NewFailErrorWithCode(
			ErrInvalidContentTypeField,
			ErrInvalidContentType,
			ErrCodeInvalidContentType,
			http.StatusUnsupportedMediaType,
		)
	}

	// Check the decoder destination
	if dest == nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			ErrNilDestination,
			gonethttp.ErrInternalServerError,
			ErrCodeNilDestination,
			http.StatusInternalServerError,
		)
	}

	// Create a new reader from the body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode JSON body into destination
	if err := decoder.Decode(dest); err != nil {
		return BodyDecodeErrorHandler(w, err)
	}
	return nil
}
