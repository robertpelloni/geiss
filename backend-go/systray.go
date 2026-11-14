package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

// SystemTrayStatus represents the runtime status of the system tray supervisor daemon.
// We provide detailed telemetry regarding platform, daemon uptime, and active state
// so the dashboard UI and external OS tray wrappers (freedesktop, Shell_NotifyIcon, NSStatusBar)
// can introspect health and control life cycle.
type SystemTrayStatus struct {
	Active     bool      `json:"active"`
	Platform   string    `json:"platform"`
	StartTime  time.Time `json:"start_time"`
	UptimeSec  int64     `json:"uptime_sec"`
	StatusText string    `json:"status_text"`
	IconState  string    `json:"icon_state"`
}

type TrayActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var (
	trayStartTime = time.Now()
	trayMu        sync.Mutex
	exitFunc      = os.Exit
	trayStatus    = SystemTrayStatus{
		Active:     true,
		Platform:   runtime.GOOS,
		StartTime:  trayStartTime,
		StatusText: "Running (Jules Autopilot System Tray Daemon)",
		IconState:  "normal",
	}
)

// systemTrayStatusHandler reports live system tray daemon health and uptime.
func systemTrayStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	trayMu.Lock()
	defer trayMu.Unlock()

	trayStatus.UptimeSec = int64(time.Since(trayStartTime).Seconds())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trayStatus)
}

// systemShutdownHandler gracefully terminates the Go server and background workers.
// Triggered either by the OS tray menu or directly from the unified web dashboard.
func systemShutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	trayMu.Lock()
	trayStatus.Active = false
	trayStatus.StatusText = "Shutting down..."
	trayStatus.IconState = "stopped"
	fn := exitFunc
	trayMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TrayActionResponse{
		Success: true,
		Message: "Graceful shutdown initiated. Process terminating in 1 second.",
	})

	// Schedule asynchronous graceful process termination to allow response delivery
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("[SYSTRAY] Graceful shutdown triggered via system tray quit endpoint.")
		if fn != nil {
			fn(0)
		}
	}()
}
