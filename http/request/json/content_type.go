package json

import (
	"net/http"
	"strings"
)

// Inspired by:
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// CheckContentType checks if the content type is JSON
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - bool: True if the content type is JSON, false otherwise
func CheckContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		mediaType := strings.ToLower(
			strings.TrimSpace(
				strings.Split(
					contentType,
					";",
				)[0],
			),
		)
		return mediaType == "application/json"
	}
	return false
}
