package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TestStepReport struct {
	Step    string `json:"step"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SystemTestReport struct {
	Timestamp string           `json:"timestamp"`
	Steps     []TestStepReport `json:"steps"`
	Overall   string           `json:"overall"`
}

func runSystemTest(repoPath string) SystemTestReport {
	report := SystemTestReport{
		Timestamp: time.Now().Format(time.RFC3339),
		Overall:   "SUCCESS",
	}

	// Step 1: Submodule Status
	_, _, errStatus := getRepoStatus(repoPath, ".") // Using root repo for a quick overall check
	if errStatus != "Clean" && errStatus != "Dirty" {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "Submodule Status Endpoint",
			Status:  "FAIL",
			Message: "Could not retrieve status",
		})
		report.Overall = "FAIL"
	} else {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "Submodule Status Endpoint",
			Status:  "OK",
			Message: "Retrieved successfully",
		})
	}

	// Step 2: Dashboard Generation
	submodules, errDash := getDashboardSubmodules(repoPath)
	if errDash != nil {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "Dashboard Generation Endpoint",
			Status:  "FAIL",
			Message: "Could not retrieve submodules",
		})
		report.Overall = "FAIL"
	} else {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "Dashboard Generation Endpoint",
			Status:  "OK",
			Message: "Submodules retrieved: " + fmt.Sprint(len(submodules)),
		})
	}

	// Step 3: Drift Detection
	driftSubs, errDrift := getSubmodulesConfig(repoPath)
	if errDrift != nil {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "Drift Detection Endpoint",
			Status:  "FAIL",
			Message: "Could not retrieve submodules",
		})
		report.Overall = "FAIL"
	} else {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "Drift Detection Endpoint",
			Status:  "OK",
			Message: "Drift scan completed for " + fmt.Sprint(len(driftSubs)) + " submodules",
		})
	}

	// Step 4: UI Auditor
	uiSubs, errUI := getSubdirectories(repoPath)
	if errUI != nil {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "UI Auditor Endpoint",
			Status:  "FAIL",
			Message: "Could not retrieve subdirectories",
		})
		report.Overall = "FAIL"
	} else {
		report.Steps = append(report.Steps, TestStepReport{
			Step:    "UI Auditor Endpoint",
			Status:  "OK",
			Message: "UI Audit completed for " + fmt.Sprint(len(uiSubs)) + " directories",
		})
	}

	return report
}

func systemTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	repoPath := ".."
	report := runSystemTest(repoPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
