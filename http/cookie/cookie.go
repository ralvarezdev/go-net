package cookie

import (
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
