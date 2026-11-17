package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ciAutoFixHandler is an honest simulation stub: no CI auto-fix command
// is configured, so it reports 501 instead of claiming success.
func ciAutoFixHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("[SHADOW] CI Auto-Fix requested (simulation: not configured)")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(AutoFixResponse{Success: false, Message: "Simulation only: no CI auto-fix command is configured; nothing was executed."})
}

type AutoFixResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
