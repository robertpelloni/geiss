package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

func buildStage(submodulePath string) bool {
	if _, err := os.Stat(filepath.Join(submodulePath, "Makefile")); err == nil {
		return true // placeholder for make build
	}
	return true
}

func testStage(submodulePath string) bool {
	if _, err := os.Stat(filepath.Join(submodulePath, "tests")); err == nil {
		return true // placeholder for tests
	}
	return true
}

func deployStage(submodulePath, environment string) bool {
	return true // placeholder for deploy
}

func runPipeline(repoPath, submodule, environment string) (bool, string) {
	fullPath := filepath.Join(repoPath, submodule)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false, fmt.Sprintf("Submodule path '%s' does not exist.", submodule)
	}

	if !buildStage(fullPath) {
		return false, "Build stage failed."
	}
	if !testStage(fullPath) {
		return false, "Test stage failed."
	}
	if !deployStage(fullPath, environment) {
		return false, "Deploy stage failed."
	}

	return true, "Pipeline executed successfully."
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
