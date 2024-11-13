package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/esa-kian/secure-guard/internal/config"
	"github.com/esa-kian/secure-guard/internal/firewall"
	"github.com/esa-kian/secure-guard/internal/monitoring"
)

func init() {
	// Set the log format and level
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

// Basic request handler function
func requestHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"ip":     r.RemoteAddr,
	}).Info("Received request")

	if firewall.CheckRequest(r) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		}).Warn("Blocked request")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - Forbidden"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to SecureGuard!"))

	duration := time.Since(startTime)
	log.WithFields(log.Fields{
		"method":   r.Method,
		"path":     r.URL.Path,
		"ip":       r.RemoteAddr,
		"duration": duration,
	}).Info("Processed request")
}

// Middleware for rate limiting
func rateLimitMiddleware(rateLimiter *firewall.RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		// Call the next handler (requestHandler) if rate limit allows
		next(w, r)
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the rate limiter with config values
	rateLimiter := firewall.NewRateLimiter(cfg.RateLimit.MaxTokens, cfg.RateLimit.RefillRate)

	// Run logging at the configured frequency
	go func() {
		ticker := time.NewTicker(cfg.Logging.Frequency)
		defer ticker.Stop()
		for range ticker.C {
			monitoring.PrintStats()
		}
	}()

	// Wrap the requestHandler with the rate limiting middleware
	http.HandleFunc("/", rateLimitMiddleware(rateLimiter, requestHandler))

	fmt.Println("Starting SecureGuard server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
