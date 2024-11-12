package firewall

import (
	"net/http"
	"strings"
)

// CheckRequest applies basic firewall rules to incoming requests.
func CheckRequest(r *http.Request) bool {
	// Example rule: Block requests from a specific user agent
	blockedUserAgents := []string{"BadBot", "Scanner"}
	for _, agent := range blockedUserAgents {
		if strings.Contains(r.UserAgent(), agent) {
			return true // Block this request
		}
	}

	// Example rule: Block requests to a specific path
	blockedPaths := []string{"/admin", "/config"}
	for _, path := range blockedPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return true // Block this request
		}
	}

	return false // Allow request if no rules match
}
