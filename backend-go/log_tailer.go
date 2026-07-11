package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

// logTailerHandler simulates streaming log content by fetching recent commit messages or dummy logs.
func logTailerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cmd := exec.Command("git", "log", "-n", "10", "--oneline")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to tail logs: %v", err)
		http.Error(w, "Failed to read logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"logs": string(output),
		"status": "success",
	})
}
