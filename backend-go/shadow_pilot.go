package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os/exec"
	"strings"
	"sync"
	"time"
)

type AnomalyReport struct {
	Timestamp time.Time `json:"timestamp"`
	Files     []string  `json:"files_modified"`
	DiffSize  int       `json:"diff_size_bytes"`
	Warning   string    `json:"warning,omitempty"`
}

func monitorGitDiff() (*AnomalyReport, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	cmd.Dir = ".."
	// Output(), not CombinedOutput(): git's CRLF "warning: in the working
	// copy..." lines go to stderr and were previously parsed as filenames.
	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("git diff failed: %v", err)
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	var modifiedFiles []string
	for _, f := range files {
		if f != "" {
			modifiedFiles = append(modifiedFiles, f)
		}
	}

	cmdFull := exec.Command("git", "diff")
	cmdFull.Dir = ".."
	fullOut, _ := cmdFull.Output()

	report := &AnomalyReport{
		Timestamp: time.Now(),
		Files:     modifiedFiles,
		DiffSize:  len(fullOut),
	}

	if len(fullOut) > 50000 {
		report.Warning = "Massive diff detected! Potential anomaly or automated mass refactoring."
	}

	return report, nil
}

// The daemon below already rescans every 30s; the handler serves its
// latest report instead of paying a full `git diff` (~4s) per request.
var (
	lastShadowReport *AnomalyReport
	shadowMu         sync.Mutex
)

func shadowPilotHandler(w http.ResponseWriter, r *http.Request) {
	shadowMu.Lock()
	report := lastShadowReport
	shadowMu.Unlock()
	if report == nil {
		var err error
		report, err = monitorGitDiff()
		if err == nil {
			lastShadowReport = report
		}
	}
	if report == nil {
		http.Error(w, "shadow scan unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func startShadowDaemon() {
	go func() {
		for {
			report, err := monitorGitDiff()
			if err != nil {
				log.Printf("[SHADOW] Error monitoring diffs: %v", err)
			} else if func() bool {
				shadowMu.Lock()
				defer shadowMu.Unlock()
				lastShadowReport = report
				return len(report.Files) > 0
			}() {
				log.Printf("[SHADOW] Detected modifications in %d files (Size: %db)", len(report.Files), report.DiffSize)
				if report.Warning != "" {
					log.Printf("[SHADOW ALERT] %s", report.Warning)
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()
}
