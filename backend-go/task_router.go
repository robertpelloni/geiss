package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var ModelCapabilities = map[string][]string{
	"Claude": {"architecture", "planning", "documentation", "refactoring", "system_understanding"},
	"Gemini": {"speed", "performance", "large_context", "scanning", "scripting"},
	"GPT":    {"code_generation", "unit_testing", "algorithm_implementation"},
}

type TaskRouteRequest struct {
	TaskDescription string `json:"task_description"`
}

type TaskRouteResponse struct {
	AssignedModel string `json:"assigned_model"`
	Reasoning     string `json:"reasoning"`
}

func analyzeTaskRequirements(taskDesc string) []string {
	taskLower := strings.ToLower(taskDesc)
	var requirements []string

	if strings.Contains(taskLower, "doc") || strings.Contains(taskLower, "plan") || strings.Contains(taskLower, "architect") || strings.Contains(taskLower, "refactor") {
		requirements = append(requirements, "architecture", "documentation", "planning")
	}

	if strings.Contains(taskLower, "scan") || strings.Contains(taskLower, "performance") || strings.Contains(taskLower, "fast") || strings.Contains(taskLower, "scripting") || strings.Contains(taskLower, "all repos") {
		requirements = append(requirements, "speed", "scripting", "large_context")
	}

	if strings.Contains(taskLower, "code") || strings.Contains(taskLower, "test") || strings.Contains(taskLower, "implement") || strings.Contains(taskLower, "algorithm") || strings.Contains(taskLower, "function") {
		requirements = append(requirements, "code_generation", "unit_testing")
	}

	return requirements
}

func routeTask(taskDesc string) (string, string) {
	if taskDesc == "" {
		return "Unknown", "No task description provided."
	}

	requirements := analyzeTaskRequirements(taskDesc)

	if len(requirements) == 0 {
		return "Claude", "Default routing applied due to ambiguous task."
	}

	scores := make(map[string]int)
	for model := range ModelCapabilities {
		scores[model] = 0
	}

	for _, req := range requirements {
		for model, capabilities := range ModelCapabilities {
			for _, cap := range capabilities {
				if req == cap {
					scores[model]++
				}
			}
		}
	}

	bestModel := ""
	maxScore := -1

	for model, score := range scores {
		if score > maxScore {
			maxScore = score
			bestModel = model
		}
	}

	if maxScore == 0 {
		return "Claude", "No strong match, defaulting to Architect."
	}

	return bestModel, "Matched requirements: " + strings.Join(requirements, ", ")
}

func taskRouterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TaskRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model, reason := routeTask(req.TaskDescription)

	resp := TaskRouteResponse{
		AssignedModel: model,
		Reasoning:     reason,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
