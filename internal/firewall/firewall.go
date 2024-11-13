package firewall

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/esa-kian/secure-guard/internal/monitoring"
)

// RateLimiter struct limits the number of requests per second
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	mu         sync.Mutex
	lastRefill time.Time
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can pass through based on rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on time elapsed since the last refill
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	refillTokens := int(elapsed / rl.refillRate)
	if refillTokens > 0 {
		rl.tokens += refillTokens
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}

	// Check if there are tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

var rateLimiter = NewRateLimiter(5, time.Second) // 5 requests per second

func CheckRequest(r *http.Request) bool {
	monitoring.RecordRequest()

	if !rateLimiter.Allow() {
		monitoring.RecordRateLimited()
		return true // Block request due to rate limit
	}

	blockedUserAgents := []string{"BadBot", "Scanner"}
	for _, agent := range blockedUserAgents {
		if strings.Contains(r.UserAgent(), agent) {
			monitoring.RecordBlockedRequest()
			return true // Block this request
		}
	}

	blockedPaths := []string{"/admin", "/config"}
	for _, path := range blockedPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			monitoring.RecordBlockedRequest()
			return true // Block this request
		}
	}

	return false // Allow request if no rules match
}
