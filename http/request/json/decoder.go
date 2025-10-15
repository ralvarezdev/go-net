package json

import (
	"encoding/json"
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Decoder struct
	Decoder struct {
		mode *goflagsmode.Flag
	}
)

// NewDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//
// Returns:
//
//   - *Decoder: The default decoder
func NewDecoder(
	mode *goflagsmode.Flag,
) *Decoder {
	return &Decoder{
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
func (d Decoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	// Check the content type
	if !CheckContentType(r) {
		return gonethttpresponse.NewFailFieldErrorWithCode(
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

	// Get the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeFailedToReadBody,
			http.StatusInternalServerError,
		)
	}

	// Decode JSON body into destination
	if err = json.Unmarshal(body, dest); err != nil {
		return BodyDecodeErrorHandler(w, err)
	}
	return nil
}
