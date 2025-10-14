package jsend

type (
	// Status is the response body status
	Status string
)

const (
	// StatusSuccess indicates a successful response
	StatusSuccess Status = "success"

	// StatusFail indicates a failed response due to client error
	StatusFail Status = "fail"

	// StatusError indicates an error response due to server error
	StatusError Status = "error"
)
