package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type PruneReport struct {
	Submodule string `json:"submodule"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	Action    string `json:"action"`
}

func getSubmodulesConfig(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "config", "--file", ".gitmodules", "--get-regexp", "path")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, nil // Return empty if no .gitmodules
	}

	var paths []string
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			paths = append(paths, parts[1])
		}
	}
	return paths, nil
}

func checkAndPruneSubmodule(repoPath, smPath string, autoPrune bool) PruneReport {
	cmdURL := exec.Command("git", "config", "--file", ".gitmodules", "--get", fmt.Sprintf("submodule.%s.url", smPath))
	cmdURL.Dir = repoPath
	outURL, err := cmdURL.Output()

	if err != nil || len(outURL) == 0 {
		return PruneReport{
			Submodule: smPath,
			Status:    "Broken (No URL configured)",
			Action:    "Requires Manual Intervention",
		}
	}

	smURL := strings.TrimSpace(string(outURL))

	// Check if remote is accessible
	cmdCheck := exec.Command("git", "ls-remote", smURL)
	cmdCheck.Dir = repoPath
	errCheck := cmdCheck.Run()

	if errCheck != nil {
		report := PruneReport{
			Submodule: smPath,
			URL:       smURL,
			Status:    "Broken (Inaccessible Remote)",
			Action:    "Flagged",
		}

		if autoPrune {
			// Remove from git cache
			cmdRm := exec.Command("git", "rm", "--cached", smPath)
			cmdRm.Dir = repoPath
			_ = cmdRm.Run()

			// Remove from .git/config
			cmdConfigRm := exec.Command("git", "config", "--remove-section", fmt.Sprintf("submodule.%s", smPath))
			cmdConfigRm.Dir = repoPath
			_ = cmdConfigRm.Run()

			// Remove from .gitmodules
			cmdGitmodRm := exec.Command("git", "config", "-f", ".gitmodules", "--remove-section", fmt.Sprintf("submodule.%s", smPath))
			cmdGitmodRm.Dir = repoPath
			_ = cmdGitmodRm.Run()

			// Remove .git/modules directory
			_ = os.RemoveAll(repoPath + "/.git/modules/" + smPath)

			report.Action = "Auto-Pruned"
		}

		return report
	}

	return PruneReport{
		Submodule: smPath,
		URL:       smURL,
		Status:    "Intact",
		Action:    "None",
	}
}

func pruneSubmodulesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	autoPrune := false
	if r.Method == http.MethodPost {
		autoPrune = true
	}

	repoPath := ".." // Root of the repository from backend-go

	submodules, err := getSubmodulesConfig(repoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get submodules: %v", err), http.StatusInternalServerError)
		return
	}

	var reports []PruneReport
	for _, sm := range submodules {
		reports = append(reports, checkAndPruneSubmodule(repoPath, sm, autoPrune))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
