package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

type LogEntry struct {
	Log string `json:"log"`
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Assuming logs are captured or redirected to a known file
	// Since go doesn't naturally keep previous log output in memory if standard log is used,
	// we would usually read a file. For this project, let's assume we tail "jules-backend.log"
	logFile := "jules-backend.log"

	file, err := os.Open(logFile)
	if err != nil {
		// Return empty if no log file is found (yet)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]LogEntry{})
		return
	}
	defer file.Close()

	var logs []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logs = append(logs, LogEntry{Log: scanner.Text()})
	}

	// For simplicity, return the last 100 lines
	if len(logs) > 100 {
		logs = logs[len(logs)-100:]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
