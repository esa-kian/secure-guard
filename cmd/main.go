package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/esa-kian/secure-guard/internal/firewall"
	"github.com/esa-kian/secure-guard/internal/monitoring"
)

// Basic request handler function
func requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	if firewall.CheckRequest(r) {
		log.Printf("Blocked request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - Forbidden"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to SecureGuard!"))
}

func main() {
	go func() {
		for {
			monitoring.PrintStats()
			time.Sleep(10 * time.Second)
		}
	}()

	http.HandleFunc("/", requestHandler)
	fmt.Println("Starting SecureGuard server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
