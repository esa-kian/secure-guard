package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/esa-kian/secure-guard/internal/firewall"
)

// Basic request handler function
func requestHandler(w http.ResponseWriter, r *http.Request) {
	// Process incoming request with firewall rules
	if firewall.CheckRequest(r) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - Forbidden"))
		return
	}

	// Allow request if it passes firewall rules
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to SecureGuard!"))
}

func main() {
	// Set up the HTTP server and route
	http.HandleFunc("/", requestHandler)

	fmt.Println("Starting SecureGuard server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
