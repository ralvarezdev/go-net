package url

import (
	"net/url"
)

// IsURL checks if a string is a valid URL
//
// Parameters:
//
//   - str: the string to check
//
// Returns:
//
//   - error: nil if the string is a valid URL, otherwise an error
func IsURL(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}
