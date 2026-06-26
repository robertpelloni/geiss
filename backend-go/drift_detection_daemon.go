package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type DriftReport struct {
	Submodule string `json:"submodule"`
	Status    string `json:"status"`
	LocalHash string `json:"local_hash,omitempty"`
	RemoteHash string `json:"remote_hash,omitempty"`
}

func getSubmodules(repoPath string) ([]string, error) {
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

func checkDrift(repoPath, smPath string) DriftReport {
	fullPath := repoPath + "/" + smPath

	if _, err := os.Stat(fullPath + "/.git"); os.IsNotExist(err) {
		return DriftReport{Submodule: smPath, Status: "Not initialized"}
	}

	// Fetch latest
	cmdFetch := exec.Command("git", "fetch", "origin")
	cmdFetch.Dir = fullPath
	_ = cmdFetch.Run()

	// Get local hash
	cmdLocal := exec.Command("git", "rev-parse", "HEAD")
	cmdLocal.Dir = fullPath
	outLocal, err := cmdLocal.Output()
	if err != nil {
		return DriftReport{Submodule: smPath, Status: "Could not determine local HEAD"}
	}
	localHash := strings.TrimSpace(string(outLocal))

	// Get remote tracking branch
	cmdBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	cmdBranch.Dir = fullPath
	outBranch, err := cmdBranch.Output()
	remoteBranch := "origin/HEAD"
	if err == nil && len(outBranch) > 0 {
		remoteBranch = strings.TrimSpace(string(outBranch))
	}

	// Get remote hash
	cmdRemote := exec.Command("git", "rev-parse", remoteBranch)
	cmdRemote.Dir = fullPath
	outRemote, err := cmdRemote.Output()
	if err != nil {
		return DriftReport{Submodule: smPath, Status: fmt.Sprintf("Could not determine remote %s", remoteBranch), LocalHash: localHash}
	}
	remoteHash := strings.TrimSpace(string(outRemote))

	if localHash != remoteHash {
		return DriftReport{
			Submodule: smPath,
			Status: "DRIFT DETECTED",
			LocalHash: localHash,
			RemoteHash: remoteHash,
		}
	}

	return DriftReport{
		Submodule: smPath,
		Status: "In sync",
		LocalHash: localHash,
		RemoteHash: remoteHash,
	}
}

func driftDetectionHandler(w http.ResponseWriter, r *http.Request) {
	repoPath := ".." // Root of the repository from backend-go

	submodules, err := getSubmodules(repoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get submodules: %v", err), http.StatusInternalServerError)
		return
	}

	var reports []DriftReport
	for _, sm := range submodules {
		reports = append(reports, checkDrift(repoPath, sm))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
