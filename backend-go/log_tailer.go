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

	logFile := "jules-backend.log"

	file, err := os.Open(logFile)
	if err != nil {
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

	if len(logs) > 100 {
		logs = logs[len(logs)-100:]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
