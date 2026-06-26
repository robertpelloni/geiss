package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type DashboardEntry struct {
	Submodule string `json:"submodule"`
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Status    string `json:"status"`
}

type DashboardResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func getDashboardSubmodules(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "config", "--file", ".gitmodules", "--get-regexp", "path")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, nil
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

func getRepoStatus(repoPath, smPath string) (string, string, string) {
	fullPath := repoPath + "/" + smPath

	if _, err := os.Stat(fullPath + "/.git"); os.IsNotExist(err) {
		return "N/A", "N/A", "Not Initialized"
	}

	cmdBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmdBranch.Dir = fullPath
	outBranch, err := cmdBranch.Output()
	branch := "Unknown"
	if err == nil {
		branch = strings.TrimSpace(string(outBranch))
	}

	cmdCommit := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmdCommit.Dir = fullPath
	outCommit, err := cmdCommit.Output()
	commit := "Unknown"
	if err == nil {
		commit = strings.TrimSpace(string(outCommit))
	}

	cmdStatus := exec.Command("git", "status", "--porcelain")
	cmdStatus.Dir = fullPath
	outStatus, _ := cmdStatus.Output()
	status := "Clean"
	if len(strings.TrimSpace(string(outStatus))) > 0 {
		status = "Dirty"
	}

	return branch, commit, status
}

func generateDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	repoPath := ".."

	submodules, err := getDashboardSubmodules(repoPath)
	if err != nil {
		http.Error(w, "Failed to get submodules", http.StatusInternalServerError)
		return
	}

	sort.Strings(submodules)

	var content strings.Builder
	content.WriteString("# Omni-Workspace Submodule Dashboard\n\n")
	content.WriteString(fmt.Sprintf("*Last updated: %s*\n\n", time.Now().Format("2006-01-02 15:04:05")))

	if len(submodules) == 0 {
		content.WriteString("No submodules detected in `.gitmodules`.\n")
	} else {
		content.WriteString("| Submodule | Branch | Commit | Status |\n")
		content.WriteString("| --- | --- | --- | --- |\n")

		for _, sm := range submodules {
			branch, commit, status := getRepoStatus(repoPath, sm)
			content.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | %s |\n", sm, branch, commit, status))
		}
	}

	err = os.WriteFile(repoPath+"/SUBMODULE_DASHBOARD.md", []byte(content.String()), 0644)
	if err != nil {
		http.Error(w, "Failed to write dashboard file", http.StatusInternalServerError)
		return
	}

	resp := DashboardResponse{
		Success: true,
		Message: "Dashboard generated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
