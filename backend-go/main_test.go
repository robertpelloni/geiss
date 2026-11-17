package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestServerStatusHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/system/server/status", nil)
	rr := httptest.NewRecorder()
	serverStatusHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp serverStatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != "online" || resp.PID == 0 {
		t.Fatalf("unexpected payload: %+v", resp)
	}
}

func TestCIAutoFixHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/shadow/autofix", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ciAutoFixHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotImplemented)
	}
	var resp AutoFixResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Success {
		t.Errorf("simulation stub must not claim success")
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

func TestDatabaseBackupHandler(t *testing.T) {
	dir := t.TempDir()
	oldDB, oldDir := backupDBPath, backupDirPath
	backupDBPath = dir + "/jules.db"
	backupDirPath = dir + "/backups"
	defer func() { backupDBPath, backupDirPath = oldDB, oldDir }()
	_ = os.WriteFile(backupDBPath, []byte("test-database-content"), 0644)

	// Test POST: Create Backup
	postReq, err := http.NewRequest("POST", "/api/system/backup", nil)
	if err != nil {
		t.Fatal(err)
	}

	postRR := httptest.NewRecorder()
	handler := http.HandlerFunc(databaseBackupHandler)
	handler.ServeHTTP(postRR, postReq)

	if status := postRR.Code; status != http.StatusOK {
		t.Errorf("POST handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var postResp BackupResponse
	if err := json.NewDecoder(postRR.Body).Decode(&postResp); err != nil {
		t.Fatal(err)
	}

	if !postResp.Success || len(postResp.Backups) == 0 {
		t.Errorf("expected backup to succeed with at least 1 backup entry, got: %+v", postResp)
	}

	// Test GET: List Backups
	getReq, err := http.NewRequest("GET", "/api/system/backup", nil)
	if err != nil {
		t.Fatal(err)
	}

	getRR := httptest.NewRecorder()
	handler.ServeHTTP(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("GET handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var getResp BackupResponse
	if err := json.NewDecoder(getRR.Body).Decode(&getResp); err != nil {
		t.Fatal(err)
	}

	if !getResp.Success || len(getResp.Backups) == 0 {
		t.Errorf("expected GET backup list to return at least 1 backup entry, got: %+v", getResp)
	}
}

func TestSystemTrayStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/system/tray", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(systemTrayStatusHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp SystemTrayStatus
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if !resp.Active {
		t.Errorf("expected tray to be active")
	}
}

func TestSystemShutdownHandler(t *testing.T) {
	trayMu.Lock()
	oldExit := exitFunc
	exitCalled := make(chan bool, 1)
	exitFunc = func(code int) {
		exitCalled <- true
	}
	trayMu.Unlock()

	req, err := http.NewRequest("POST", "/api/system/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(systemShutdownHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp TrayActionResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if !resp.Success {
		t.Errorf("expected success to be true")
	}

	select {
	case <-exitCalled:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Errorf("expected exitFunc to be called")
	}

	trayMu.Lock()
	exitFunc = oldExit
	trayMu.Unlock()
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

func TestAuthMiddlewareFailClosed(t *testing.T) {
	ok := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }
	h := authMiddleware(ok)
	t.Setenv("JULES_API_KEY", "")
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodPost, "/x", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unset key: want 401 got %d", rr.Code)
	}
	t.Setenv("JULES_API_KEY", "secret")
	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	req.Header.Set("X-API-KEY", "wrong")
	h(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("wrong key: want 401 got %d", rr.Code)
	}
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/x", nil)
	req.Header.Set("X-API-KEY", "secret")
	h(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("good key: want 200 got %d", rr.Code)
	}
}

func TestTelemetryTraversalRejected(t *testing.T) {
	b, _ := json.Marshal(TelemetryRequest{TargetDir: "../../outside"})
	req := httptest.NewRequest(http.MethodPost, "/api/system/telemetry", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	telemetryStandardizerHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("traversal: want 400 got %d", rr.Code)
	}
}

func TestStateSnapshotNoDB(t *testing.T) {
	dir := t.TempDir()
	info, err := createStateSnapshot(dir + "/backups")
	if err != nil {
		t.Fatal(err)
	}
	if info == nil || info.Filename == "" {
		t.Fatalf("bad snapshot: %+v", info)
	}
	if _, err := os.Stat(dir + "/backups/" + info.Filename); err != nil {
		t.Fatal(err)
	}
}

func TestPipelineNonexistentSubmoduleFails(t *testing.T) {
	ok, _ := runPipeline("..", "definitely_missing_submodule_xyz", "local")
	if ok {
		t.Fatal("nonexistent submodule must fail, not succeed")
	}
}
