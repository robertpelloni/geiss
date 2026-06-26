package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var backendExts = []string{".py", ".go", ".rs", ".java", ".c", ".cpp", ".rb"}
var frontendExts = []string{".html", ".jsx", ".tsx", ".vue", ".css", ".scss", ".js", ".ts"}

type UIAuditReport struct {
	Submodule     string `json:"submodule"`
	BackendFiles  int    `json:"backend_files"`
	FrontendFiles int    `json:"frontend_files"`
	Status        string `json:"status"`
}

func hasExtension(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func scanDirectoryForExtensions(directory string, extensions []string) int {
	count := 0
	_ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "venv" {
				return filepath.SkipDir
			}
			return nil
		}

		if hasExtension(info.Name(), extensions) {
			count++
		}
		return nil
	})
	return count
}

func auditSubmodule(repoPath, smPath string) UIAuditReport {
	fullPath := filepath.Join(repoPath, smPath)

	backendCount := scanDirectoryForExtensions(fullPath, backendExts)
	frontendCount := scanDirectoryForExtensions(fullPath, frontendExts)

	status := "OK"
	if backendCount > 0 && frontendCount == 0 {
		status = "WARNING: Backend logic detected with ZERO frontend files."
	} else if backendCount > 10 && frontendCount < 2 {
		status = "WARNING: High backend-to-frontend ratio. UI may be incomplete."
	}

	return UIAuditReport{
		Submodule:     smPath,
		BackendFiles:  backendCount,
		FrontendFiles: frontendCount,
		Status:        status,
	}
}

func getSubdirectories(repoPath string) ([]string, error) {
	var dirs []string
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			name := entry.Name()
			if name != "scripts" && name != "docs" && name != "tests" && name != "backend-go" && name != "src" && name != "dist" {
				dirs = append(dirs, name)
			}
		}
	}
	return dirs, nil
}

func uiAuditorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	repoPath := ".."

	submodules, err := getSubdirectories(repoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get subdirectories: %v", err), http.StatusInternalServerError)
		return
	}

	var reports []UIAuditReport
	for _, sm := range submodules {
		reports = append(reports, auditSubmodule(repoPath, sm))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
