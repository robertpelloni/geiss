package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ConflictBlock struct {
	Local  []string `json:"local"`
	Remote []string `json:"remote"`
}

type ConflictReport struct {
	File     string          `json:"file"`
	Blocks   []ConflictBlock `json:"blocks"`
	Strategy []string        `json:"strategy"`
}

func checkForConflicts(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var files []string
	for _, f := range lines {
		if f != "" {
			files = append(files, f)
		}
	}
	return files, nil
}

func extractConflictBlocks(filePath string) ([]ConflictBlock, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var conflicts []ConflictBlock
	var currentConflict *ConflictBlock
	state := "normal" // normal, local, remote

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "<<<<<<< ") {
			currentConflict = &ConflictBlock{}
			state = "local"
		} else if strings.HasPrefix(line, "=======") {
			if state == "local" {
				state = "remote"
			}
		} else if strings.HasPrefix(line, ">>>>>>> ") {
			if currentConflict != nil && state == "remote" {
				conflicts = append(conflicts, *currentConflict)
				currentConflict = nil
				state = "normal"
			}
		} else {
			if state == "local" && currentConflict != nil {
				currentConflict.Local = append(currentConflict.Local, line)
			} else if state == "remote" && currentConflict != nil {
				currentConflict.Remote = append(currentConflict.Remote, line)
			}
		}
	}

	return conflicts, scanner.Err()
}

func analyzeConflict(block ConflictBlock) string {
	localContent := strings.TrimSpace(strings.Join(block.Local, ""))
	remoteContent := strings.TrimSpace(strings.Join(block.Remote, ""))

	if localContent == "" && remoteContent != "" {
		return "Accept Remote (Local is empty)"
	}
	if localContent != "" && remoteContent == "" {
		return "Accept Local (Remote is empty)"
	}

	return "Requires Intelligent Merge (Both contain features)"
}

func conflictResolutionHandler(w http.ResponseWriter, r *http.Request) {
	repoPath := ".." // Root of the repository from backend-go

	files, err := checkForConflicts(repoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check for conflicts: %v", err), http.StatusInternalServerError)
		return
	}

	var reports []ConflictReport
	for _, file := range files {
		fullPath := repoPath + "/" + file
		blocks, err := extractConflictBlocks(fullPath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var strategies []string
		for _, block := range blocks {
			strategies = append(strategies, analyzeConflict(block))
		}

		reports = append(reports, ConflictReport{
			File:     file,
			Blocks:   blocks,
			Strategy: strategies,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
