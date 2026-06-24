package main

import (
	"encoding/json"
	"net/http"
)

func queueTelemetryHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	db.Order("created_at desc").Limit(50).Find(&tasks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
