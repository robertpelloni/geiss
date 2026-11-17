package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
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
	statuses := []SubmoduleStatus{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hashPart := parts[0]
			statusFlag := "Synchronized"

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

			statuses = append(statuses, SubmoduleStatus{
				Commit: hashPart,
				Path:   parts[1],
				Status: statusFlag,
			})
		}
	}

	return statuses, nil
}

// Why: `git submodule status` takes ~4-9s on this tree, which made the
// dashboard panel look broken. Cache the last good scan for 30s; the
// data changes only on git operations, never faster than that matters.
var statusCache struct {
	sync.Mutex
	at   time.Time
	data []SubmoduleStatus
}

func cachedSubmodules() ([]SubmoduleStatus, error) {
	statusCache.Lock()
	defer statusCache.Unlock()
	if statusCache.data != nil && time.Since(statusCache.at) < 30*time.Second {
		return statusCache.data, nil
	}
	statuses, err := checkSubmodules()
	if err != nil {
		return nil, err
	}
	statusCache.data, statusCache.at = statuses, time.Now()
	return statuses, nil
}

func submoduleStatusHandler(w http.ResponseWriter, r *http.Request) {
	statuses, err := cachedSubmodules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}
