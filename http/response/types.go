package response

type (
	// ErrorResponse struct
	ErrorResponse struct {
		Error string `json:"error"`
	}
)

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}
