package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
		initDB()
	// Register API handlers
	http.HandleFunc("/api/shadow/diff", shadowPilotHandler)
	http.HandleFunc("/system/status", submoduleStatusHandler)
	http.HandleFunc("/api/shadow/autofix", authMiddleware(ciAutoFixHandler))
	http.HandleFunc("/api/queue/telemetry", queueTelemetryHandler)
	http.HandleFunc("/api/tasks/route", taskRouterHandler)
	http.HandleFunc("/api/conflicts/resolve", authMiddleware(conflictResolutionHandler))
	http.HandleFunc("/api/system/drift", driftDetectionHandler)
	http.HandleFunc("/api/system/prune", authMiddleware(pruneSubmodulesHandler))
	http.HandleFunc("/api/system/dashboard", generateDashboardHandler)
	http.HandleFunc("/api/system/audit", uiAuditorHandler)
	http.HandleFunc("/api/system/telemetry", authMiddleware(telemetryStandardizerHandler))
	http.HandleFunc("/api/system/test", systemTestHandler)

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
