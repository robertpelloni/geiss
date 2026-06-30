package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SearchReplaceRequest struct {
	SearchStr  string `json:"search_str"`
	ReplaceStr string `json:"replace_str"`
	DryRun     bool   `json:"dry_run"`
}

type SearchReplaceReport struct {
	ModifiedFiles int      `json:"modified_files"`
	Files         []string `json:"files"`
}

func processFile(filePath, searchStr, replaceStr string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	contentStr := string(content)
	if strings.Contains(contentStr, searchStr) {
		newContent := strings.ReplaceAll(contentStr, searchStr, replaceStr)
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		return err == nil
	}
	return false
}

func searchAndReplaceInDir(directory, searchStr, replaceStr string, dryRun bool) ([]string, int) {
	var modifiedFiles []string
	var count int

	_ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "venv" || name == "dist" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		if dryRun {
			content, err := os.ReadFile(path)
			if err == nil && strings.Contains(string(content), searchStr) {
				modifiedFiles = append(modifiedFiles, path)
				count++
			}
		} else {
			if processFile(path, searchStr, replaceStr) {
				modifiedFiles = append(modifiedFiles, path)
				count++
			}
		}

		return nil
	})

	return modifiedFiles, count
}

func globalSearchAndReplaceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SearchReplaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.SearchStr == "" {
		http.Error(w, "Search string cannot be empty", http.StatusBadRequest)
		return
	}

	repoPath := ".."

	submodules, err := getDashboardSubmodules(repoPath)
	if err != nil {
		http.Error(w, "Failed to get submodules", http.StatusInternalServerError)
		return
	}

	report := SearchReplaceReport{
		Files: make([]string, 0),
	}

	for _, sm := range submodules {
		fullPath := filepath.Join(repoPath, sm)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			files, count := searchAndReplaceInDir(fullPath, req.SearchStr, req.ReplaceStr, req.DryRun)
			report.ModifiedFiles += count
			report.Files = append(report.Files, files...)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
