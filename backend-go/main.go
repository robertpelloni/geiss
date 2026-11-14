package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// No database: the task queue is in-memory (queue.go). The former
	// gorm/sqlite layer was write-only dead code and broke the build
	// on current toolchains; it is archived as db.go.bak.
	// Register API handlers
	http.HandleFunc("/api/shadow/diff", shadowPilotHandler)
	http.HandleFunc("/system/status", submoduleStatusHandler)
	http.HandleFunc("/api/shadow/autofix", ciAutoFixHandler)
	http.HandleFunc("/api/queue/telemetry", queueTelemetryHandler)
	http.HandleFunc("/api/tasks/route", taskRouterHandler)
	http.HandleFunc("/api/conflicts/resolve", conflictResolutionHandler)
	http.HandleFunc("/api/system/drift", driftDetectionHandler)
	http.HandleFunc("/api/system/prune", pruneSubmodulesHandler)
	http.HandleFunc("/api/system/dashboard", authMiddleware(generateDashboardHandler))
	http.HandleFunc("/api/system/audit", authMiddleware(uiAuditorHandler))
	http.HandleFunc("/api/system/telemetry", telemetryStandardizerHandler)
	http.HandleFunc("/api/system/logs", logTailerHandler)
	http.HandleFunc("/api/system/pipeline", deploymentPipelineHandler)
	http.HandleFunc("/api/system/refactor", globalSearchAndReplaceHandler)
	http.HandleFunc("/api/system/server/status", serverStatusHandler)
	http.HandleFunc("/api/system/server/shutdown", authMiddleware(serverShutdownHandler))

	// Start background daemons
	startShadowDaemon()
	StartQueueWorker()

	// Serve the static SPA build
	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)

	fmt.Println("Jules Autopilot Go Backend starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
