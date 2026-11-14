package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// Why: the system tray (scripts/system_tray.py) and the dashboard both poll
// status and request graceful shutdown. These handlers previously existed
// only as client calls with no server side, so the tray could never quit
// the backend. startTime is captured once at process start.
var startTime = time.Now()

type serverStatusResponse struct {
	Status    string `json:"status"`
	PID       int    `json:"pid"`
	UptimeSec int64  `json:"uptime_sec"`
	StartedAt string `json:"started_at"`
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(serverStatusResponse{
		Status:    "online",
		PID:       os.Getpid(),
		UptimeSec: int64(time.Since(startTime).Seconds()),
		StartedAt: startTime.Format(time.RFC3339),
	})
}

func serverShutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "shutdown initiated"})
	// Flush before exiting so the tray/dashboard sees the confirmation.
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	go func() {
		time.Sleep(200 * time.Millisecond)
		os.Exit(0)
	}()
}
