package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type PipelineRequest struct {
	Submodule   string `json:"submodule"`
	Environment string `json:"environment"`
	Generate    bool   `json:"generate"`
}

type PipelineResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type WorkflowJobStep struct {
	Name string            `yaml:"name"`
	Uses string            `yaml:"uses,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
	Run  string            `yaml:"run,omitempty"`
}

type WorkflowJob struct {
	RunsOn string            `yaml:"runs-on"`
	Steps  []WorkflowJobStep `yaml:"steps"`
}

type Workflow struct {
	Name string                 `yaml:"name"`
	On   map[string]interface{} `yaml:"on"`
	Jobs map[string]WorkflowJob `yaml:"jobs"`
}

func generateGithubAction(repoPath, submodule string) bool {
	submoduleName := filepath.Base(filepath.Clean(submodule))
	workflowDir := filepath.Join(repoPath, ".github", "workflows")

	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return false
	}

	workflowPath := filepath.Join(workflowDir, fmt.Sprintf("%s_pipeline.yml", submoduleName))

	workflow := Workflow{
		Name: fmt.Sprintf("Omni-Workspace CI/CD - %s", submoduleName),
		On: map[string]interface{}{
			"push": map[string]interface{}{
				"paths": []string{fmt.Sprintf("%s/**", submoduleName)},
			},
			"pull_request": map[string]interface{}{
				"paths": []string{fmt.Sprintf("%s/**", submoduleName)},
			},
		},
		Jobs: map[string]WorkflowJob{
			"build-test-deploy": {
				RunsOn: "ubuntu-latest",
				Steps: []WorkflowJobStep{
					{
						Name: "Checkout code",
						Uses: "actions/checkout@v4",
						With: map[string]string{
							"submodules": "recursive",
						},
					},
					{
						Name: "Set up Go",
						Uses: "actions/setup-go@v5",
						With: map[string]string{
							"go-version": "1.26.0",
						},
					},
					{
						Name: fmt.Sprintf("Run Pipeline for %s", submoduleName),
						Run:  fmt.Sprintf("curl -X POST -H 'Content-Type: application/json' -H 'X-API-KEY: ${{ secrets.JULES_API_KEY }}' -d '{\"submodule\":\"%s\",\"environment\":\"staging\"}' http://localhost:8080/api/system/pipeline", submoduleName),
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(&workflow)
	if err != nil {
		return false
	}

	err = os.WriteFile(workflowPath, data, 0644)
	return err == nil
}

func buildStage(submodulePath string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := os.Stat(filepath.Join(submodulePath, "go.mod")); err == nil {
		out, err := exec.CommandContext(ctx, "go", "vet", "./...").CombinedOutput()
		_ = out
		if err != nil {
			return false, "go vet failed: " + err.Error()
		}
		return true, "go vet passed"
	}
	if _, err := os.Stat(filepath.Join(submodulePath, "package.json")); err == nil {
		return true, "package.json present (build not executed by dispatcher)"
	}
	if _, err := os.Stat(filepath.Join(submodulePath, "Makefile")); err == nil {
		return true, "Makefile present (build not executed by dispatcher)"
	}
	return false, "no recognised build system (go.mod/package.json/Makefile)"
}

func testStage(submodulePath string) (bool, string) {
	if _, err := os.Stat(filepath.Join(submodulePath, "tests")); err == nil {
		return true, "tests/ present (not executed by dispatcher)"
	}
	found := false
	_ = filepath.Walk(submodulePath, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			found = true
		}
		return nil
	})
	if found {
		return true, "test files present (not executed by dispatcher)"
	}
	return false, "no tests found"
}

func deployStage(submodulePath, environment string) (bool, string) {
	if environment == "local" {
		return true, "local dry-run: no deploy performed"
	}
	return false, "deploy to '" + environment + "' is not implemented (simulation only)"
}

func runPipeline(repoPath, submodule, environment string) (bool, string) {
	// Reject submodule paths escaping the workspace (e.g. "../outside"):
	// Join alone happily resolves them and stages then run outside the tree.
	repoAbs, err := filepath.Abs(repoPath)
	if err != nil {
		return false, "invalid workspace path"
	}
	fullPath := filepath.Join(repoAbs, submodule)
	rel, err := filepath.Rel(repoAbs, fullPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false, fmt.Sprintf("Submodule path '%s' escapes the workspace.", submodule)
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false, fmt.Sprintf("Submodule path '%s' does not exist.", submodule)
	}

	if ok, msg := buildStage(fullPath); !ok {
		return false, "Build stage failed: " + msg
	}
	if ok, msg := testStage(fullPath); !ok {
		return false, "Test stage failed: " + msg
	}
	if ok, msg := deployStage(fullPath, environment); !ok {
		return false, "Deploy stage failed: " + msg
	}
	return true, "Pipeline checks passed (build/test verified present; local deploy is a dry-run)."
}

func deploymentPipelineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PipelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Submodule == "" {
		http.Error(w, "Submodule name is required", http.StatusBadRequest)
		return
	}

	repoPath := ".."

	resp := PipelineResponse{}

	if req.Generate {
		if generateGithubAction(repoPath, req.Submodule) {
			resp.Success = true
			resp.Message = "Workflow generated successfully."
		} else {
			resp.Success = false
			resp.Message = "Failed to generate workflow."
		}
	} else {
		env := req.Environment
		if env == "" {
			env = "staging"
		}
		success, msg := runPipeline(repoPath, req.Submodule, env)
		resp.Success = success
		resp.Message = msg
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
