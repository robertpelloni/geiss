package main

import (
	"encoding/json"
	"net/http"
)

func queueTelemetryHandler(w http.ResponseWriter, r *http.Request) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queue.Tasks)
}
