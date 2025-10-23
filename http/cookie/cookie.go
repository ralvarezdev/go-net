package cookie

import (
	"net/http"
	"time"
)

type (
	// NotFoundFn is the function to call when a cookie is not found
	NotFoundFn func(
		w http.ResponseWriter,
		attributes *Attributes,
		err error,
	) error

	// Attributes is the structure for the attributes of a cookie
	Attributes struct {
		Name     string
		Path     string
		Domain   string
		Secure   bool
		HTTPOnly bool
		SameSite http.SameSite
	}
)

// SetCookie sets a cookie
//
// Parameters:
//
//   - w: The HTTP response writer
//   - attributes: The attributes of the cookie
//   - value: The value of the cookie
//   - expiresAt: The expiration time of the cookie
func SetCookie(
	w http.ResponseWriter,
	attributes *Attributes,
	value string,
	expiresAt time.Time,
) {
	// Create and create cookie
	cookie := &http.Cookie{
		Name:     attributes.Name,
		Value:    value,
		Path:     attributes.Path,
		Domain:   attributes.Domain,
		Expires:  expiresAt,
		Secure:   attributes.Secure,
		HttpOnly: attributes.HTTPOnly,
		SameSite: attributes.SameSite,
	}
	http.SetCookie(w, cookie)
}

// SetTimestampCookie sets a cookie with a timestamp
//
// Parameters:
//
//   - w: The HTTP response writer
//   - attributes: The attributes of the cookie
//   - value: The value of the cookie
//   - expiresAt: The expiration time of the cookie
func SetTimestampCookie(
	w http.ResponseWriter,
	attributes *Attributes,
	value,
	expiresAt time.Time,
) {
	SetCookie(
		w,
		attributes,
		value.Format(time.RFC3339),
		expiresAt,
	)
}

// GetTimestampCookie gets a timestamp cookie
//
// Parameters:
//
//   - r: The HTTP request
//   - attributes: The attributes of the cookie
//
// Returns:
//
//   - *time.Time: The value of the cookie, or nil if not found
//   - error: An error if something went wrong
func GetTimestampCookie(
	r *http.Request,
	attributes *Attributes,
) (*time.Time, error) {
	// Check if the request or the attributes is nil
	if r == nil {
		return nil, ErrNilRequest
	}
	if attributes == nil {
		return nil, ErrNilAttributes
	}

	// Get the cookie
	cookie, err := r.Cookie(attributes.Name)
	if err != nil {
		return nil, err
	}

	// Parse the cookie value
	value, err := time.Parse(time.RFC3339, cookie.Value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// DeleteCookies deletes the cookies
//
// Parameters:
//
//   - w: The HTTP response writer
//   - attributesList: The list of attributes of the cookies to delete
func DeleteCookies(w http.ResponseWriter, attributesList ...*Attributes) {
	for _, attributes := range attributesList {
		SetCookie(
			w,
			attributes,
			"",
			time.Now().Add(-time.Hour),
		)
	}
}

// RenovateCookie creates a new cookie with the same value and a new expiration time
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - attributes: The attributes of the cookie
//   - expiresAt: The new expiration time of the cookie
//   - cookieNotFoundFn: The function to call if the cookie is not found
//   - error: An error if something went wrong
//
// Returns:
//
//   - error: An error if something went wrong
func RenovateCookie(
	w http.ResponseWriter,
	r *http.Request,
	attributes *Attributes,
	expiresAt time.Time,
	cookieNotFoundFn NotFoundFn,
) error {
	// Check if the request or the attributes is nil
	if r == nil {
		return ErrNilRequest
	}
	if attributes == nil {
		return ErrNilAttributes
	}

	// Get the cookie
	cookie, err := r.Cookie(attributes.Name)
	if err != nil {
		return cookieNotFoundFn(w, attributes, err)
	}
	SetCookie(
		w,
		attributes,
		cookie.Value,
		expiresAt,
	)
	return nil
}
