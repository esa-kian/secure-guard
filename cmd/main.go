package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

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

func main() {
	go func() {
		for {
			monitoring.PrintStats()
			time.Sleep(10 * time.Second)
		}
	}()

	http.HandleFunc("/", requestHandler)
	log.Info("Starting SecureGuard server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
