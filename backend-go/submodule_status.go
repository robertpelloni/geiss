package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type SubmoduleStatus struct {
	Path   string `json:"path"`
	Commit string `json:"commit"`
	Status string `json:"status"`
}

func checkSubmodules() ([]SubmoduleStatus, error) {
	cmd := exec.Command("git", "submodule", "status")
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git submodule status failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var statuses []SubmoduleStatus
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
	for _, line := range lines {
		if line == "" {
			continue
		}
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hashPart := parts[0]
			statusFlag := "Synchronized"
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
			if strings.HasPrefix(hashPart, "-") {
				statusFlag = "Not Initialized"
				hashPart = hashPart[1:]
			} else if strings.HasPrefix(hashPart, "+") {
				statusFlag = "Out of Sync (Detached/Modified)"
				hashPart = hashPart[1:]
			} else if strings.HasPrefix(hashPart, "U") {
				statusFlag = "Merge Conflict"
				hashPart = hashPart[1:]
			}
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
			statuses = append(statuses, SubmoduleStatus{
				Commit: hashPart,
				Path:   parts[1],
				Status: statusFlag,
			})
		}
	}
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
	return statuses, nil
}

func submoduleStatusHandler(w http.ResponseWriter, r *http.Request) {
	statuses, err := checkSubmodules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}
