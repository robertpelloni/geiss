package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	// Setup mock DB for testing

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
