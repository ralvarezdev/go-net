package http

import (
	"github.com/ralvarezdev/go-net/http/status"
	"net/http"
	"strings"
)

// GetClientIP returns the client's IP address from the request
func GetClientIP(r *http.Request) string {
	// Check if the request has a forwarded IP from a proxy or load balancer
	forwarded := r.Header.Get(status.XForwardedFor)
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IP addresses, the client's IP is the first one
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip)
	}

	// If there's no forwarded IP, use RemoteAddr
	ip := r.RemoteAddr
	if len(ip) > 0 && ip[0] == '[' {
		// Handle IPv6 address
		ip = strings.Split(ip, "]")[0] + "]"
	} else {
		// Handle IPv4 address
		ip = strings.Split(ip, ":")[0]
	}
	return strings.TrimSpace(ip)
}
