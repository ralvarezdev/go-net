package json

import (
	"encoding/json"
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
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

// DecodeReader decodes the JSON body and stores it in the destination
//
// Parameters:
//
//   - reader: The reader to read the body from
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeReader(
	reader io.Reader,
	dest interface{},
) (err error) {
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
	body, err := io.ReadAll(reader)
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
		return BodyDecodeErrorHandler(err)
	}
	return nil
}

// DecodeRequest decodes the JSON request body and stores it in the destination
//
// Parameters:
//
//   - request: The HTTP request
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeRequest(
	request *http.Request,
	dest interface{},
) error {
	// Check the request
	if request == nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrNilRequest,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeNilRequest,
			http.StatusInternalServerError,
		)
	}

	// Check the content type
	if !gonethttprequest.CheckContentType(request) {
		return gonethttpresponse.NewFailFieldErrorWithCode(
			gonethttprequest.ErrInvalidContentTypeField,
			gonethttprequest.ErrInvalidContentType,
			gonethttprequest.ErrCodeInvalidContentType,
			http.StatusUnsupportedMediaType,
		)
	}
	return d.DecodeReader(request.Body, dest)
}
