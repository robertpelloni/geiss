package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Set up file logging
	logFile, err := os.OpenFile("jules-backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
	} else {
		log.Println("Failed to open log file, using default stderr")
	}

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
	http.HandleFunc("/api/system/dashboard", authMiddleware(generateDashboardHandler))
	http.HandleFunc("/api/system/audit", authMiddleware(uiAuditorHandler))
	http.HandleFunc("/api/system/telemetry", authMiddleware(telemetryStandardizerHandler))
	http.HandleFunc("/api/system/test", authMiddleware(systemTestHandler))
	http.HandleFunc("/api/system/refactor", authMiddleware(globalSearchAndReplaceHandler))
	http.HandleFunc("/api/system/pipeline", authMiddleware(deploymentPipelineHandler))
	http.HandleFunc("/api/system/logs", authMiddleware(getLogsHandler))

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
