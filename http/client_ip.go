package http

import (
	"net/http"
	"strings"
)

// GetClientIP returns the client's IP address from the request
//
// It checks the X-Forwarded-For header first (in case of proxies or load balancers),
// and falls back to the RemoteAddr if the header is not present.
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - string: The client's IP address
func GetClientIP(r *http.Request) string {
	// Check if the request has a forwarded IP from a proxy or load balancer
	forwarded := r.Header.Get(XForwardedFor)
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IP addresses, the client's IP is the first one
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip)
	}

	// If there's no forwarded IP, use RemoteAddr
	ip := r.RemoteAddr
	if ip != "" && ip[0] == '[' {
		// Handle IPv6 address
		ip = strings.Split(ip, "]")[0] + "]"
	} else {
		// Handle IPv4 address
		ip = strings.Split(ip, ":")[0]
	}
	return strings.TrimSpace(ip)
}
