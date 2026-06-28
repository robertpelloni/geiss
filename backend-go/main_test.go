package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func setupTestDB() {
	// Use an in-memory sqlite db for tests
	_ = os.Remove("jules_test.db")
	db, _ = gorm.Open(sqlite.Open("jules_test.db"), &gorm.Config{})
	db.AutoMigrate(&Task{})
}

func teardownTestDB() {
	_ = os.Remove("jules_test.db")
}

func TestCIAutoFixHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/shadow/autofix", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ciAutoFixHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp AutoFixResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if !resp.Success {
		t.Errorf("expected success to be true")
	}
}

func TestTaskRouterHandler(t *testing.T) {
	payload := TaskRouteRequest{TaskDescription: "Refactor the backend structure"}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/api/tasks/route", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskRouterHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp TaskRouteResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.AssignedModel != "Claude" {
		t.Errorf("expected Claude, got %s", resp.AssignedModel)
	}
}

func TestQueueTelemetryHandler(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	req, err := http.NewRequest("GET", "/api/queue/telemetry", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(queueTelemetryHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestConflictResolutionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/conflicts/resolve", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(conflictResolutionHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp []ConflictReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	// Given a fresh start, there should be no conflicts
	if len(resp) != 0 {
		t.Errorf("expected 0 conflicts, got %d", len(resp))
	}
}

func TestDriftDetectionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/system/drift", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(driftDetectionHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp []DriftReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestPruneSubmodulesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/system/prune", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pruneSubmodulesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp []PruneReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateDashboardHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/system/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(generateDashboardHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp DashboardResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if !resp.Success {
		t.Errorf("expected success to be true")
	}
}

func TestSubmoduleStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/system/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(submoduleStatusHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp []SubmoduleStatus
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestUIAuditorHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/system/audit", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uiAuditorHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp []UIAuditReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestTelemetryStandardizerHandler(t *testing.T) {
	payload := TelemetryRequest{TargetDir: "."}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/api/system/telemetry", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(telemetryStandardizerHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp []TelemetryReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestGlobalSearchAndReplaceHandler(t *testing.T) {
	payload := SearchReplaceRequest{
		SearchStr:  "test_string",
		ReplaceStr: "new_string",
		DryRun:     true,
	}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/api/system/refactor", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(globalSearchAndReplaceHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp SearchReplaceReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestDeploymentPipelineHandler(t *testing.T) {
	payload := PipelineRequest{
		Submodule:   "dummy_submodule",
		Environment: "staging",
		Generate:    false,
	}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/api/system/pipeline", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deploymentPipelineHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp PipelineResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestSystemTestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/system/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(systemTestHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp SystemTestReport
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestGetLogsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/system/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getLogsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
