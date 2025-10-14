package jsend

/*
// ParseGRPCError parses a gRPC error to a JSend error response
//
// Parameters:
//
//   - mode: the flag mode
//   - err: the gRPC error
//
// Returns:
//   - *gonethttpresponsejsend.ErrorResponse: the JSend error response
func ParseGRPCError(
	mode *goflagsmode.Flag,
	err error,
) *gonethttpresponsejsend.ErrorResponse {
	// Check if the mode or the error is nil
	if mode == nil {
		return gonethttpresponsejsend.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("nil mode"),
			nil,
			mode,
		)
	}
	if err == nil {
		return gonethttpresponsejsend.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("nil error"),
			nil,
			mode,
		)
	}

	// Parse the gRPC error to a JSend error response
	return gonethttpresponsejsend.ParseGRPCError(err, mode)
}
*/
