package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Standard Omni-Workspace Telemetry Stubs
const pythonLogStub = `
import logging
# Omni-Workspace Standard Telemetry
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
if not logger.handlers:
    handler = logging.StreamHandler()
    formatter = logging.Formatter('%(asctime)s - [OMNI] - %(name)s - %(levelname)s - %(message)s')
    handler.setFormatter(formatter)
    logger.addHandler(handler)
`

const nodejsLogStub = `
// Omni-Workspace Standard Telemetry
const omniLogger = {
    info: (msg) => console.log(` + "`" + `[${new Date().toISOString()}] - [OMNI] - INFO - ${msg}` + "`" + `),
    error: (msg) => console.error(` + "`" + `[${new Date().toISOString()}] - [OMNI] - ERROR - ${msg}` + "`" + `)
};
`

type TelemetryReport struct {
	File    string `json:"file"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TelemetryRequest struct {
	TargetDir string `json:"target_dir"`
}

func detectLanguage(filePath string) string {
	if strings.HasSuffix(filePath, ".py") {
		return "python"
	}
	if strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".ts") {
		return "nodejs"
	}
	return "unknown"
}

func checkFileForTelemetry(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "[OMNI]")
}

func standardizeFile(filePath string) (bool, string) {
	if checkFileForTelemetry(filePath) {
		return false, "Telemetry already present"
	}

	lang := detectLanguage(filePath)
	var stub string

	switch lang {
	case "python":
		stub = pythonLogStub
	case "nodejs":
		stub = nodejsLogStub
	default:
		return false, "Unsupported language for " + filePath
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, "Error reading file: " + err.Error()
	}

	lines := strings.Split(string(content), "\n")
	var newContent string

	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!") {
		newContent = lines[0] + "\n" + stub + "\n" + strings.Join(lines[1:], "\n")
	} else {
		newContent = stub + "\n" + string(content)
	}

	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return false, "Error writing to file: " + err.Error()
	}

	return true, "Telemetry injected successfully"
}

func scanAndStandardize(targetDir string) []TelemetryReport {
	var results []TelemetryReport

	_ = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "venv" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		lang := detectLanguage(path)
		if lang != "unknown" {
			success, msg := standardizeFile(path)
			results = append(results, TelemetryReport{
				File:    path,
				Success: success,
				Message: msg,
			})
		}
		return nil
	})

	return results
}

func telemetryStandardizerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TelemetryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.TargetDir == "" {
		req.TargetDir = ".."
	} else {
		req.TargetDir = filepath.Join("..", req.TargetDir)
	}

    // Path traversal mitigation
	absPath, err := filepath.Abs(req.TargetDir)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

    repoRoot, _ := filepath.Abs("..")
	repoRoot = repoRoot + string(filepath.Separator)
	if !strings.HasPrefix(absPath, repoRoot) {
		http.Error(w, "Access denied: Target directory is outside of the repository", http.StatusForbidden)
		return
	}

	reports := scanAndStandardize(req.TargetDir)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
