package cookie

import (
	"errors"
	"net/http"
	"time"
)

// Attributes is the structure for the attributes of a cookie
type Attributes struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// SetCookie sets a cookie
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
func SetTimestampCookie(
	w http.ResponseWriter,
	attributes *Attributes,
	expiresAt time.Time,
) {
	SetCookie(
		w,
		attributes,
		expiresAt.Format(time.RFC3339),
		expiresAt,
	)
}

// GetTimestampCookie gets a timestamp cookie
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

	// Check if the error is not found
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
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
func RenovateCookie(
	w http.ResponseWriter,
	r *http.Request,
	attributes *Attributes,
	expiresAt time.Time,
	cookieNotFoundFn func(
		w http.ResponseWriter,
		attributes *Attributes,
		err error,
	) error,
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
