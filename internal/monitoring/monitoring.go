package monitoring

import (
	"fmt"
	"sync"
)

type Monitor struct {
	TotalRequests   int
	BlockedRequests int
	RateLimited     int
	mu              sync.Mutex
}

var monitorInstance = &Monitor{}

func RecordRequest() {
	monitorInstance.mu.Lock()
	monitorInstance.TotalRequests++
	monitorInstance.mu.Unlock()
}

func RecordBlockedRequest() {
	monitorInstance.mu.Lock()
	monitorInstance.BlockedRequests++
	monitorInstance.mu.Unlock()
}

func RecordRateLimited() {
	monitorInstance.mu.Lock()
	monitorInstance.RateLimited++
	monitorInstance.mu.Unlock()
}

func PrintStats() {
	monitorInstance.mu.Lock()
	defer monitorInstance.mu.Unlock()
	fmt.Printf("Total Requests: %d, Blocked: %d, Rate Limited: %d\n",
		monitorInstance.TotalRequests, monitorInstance.BlockedRequests, monitorInstance.RateLimited)
}
