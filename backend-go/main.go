package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Register API handlers
	http.HandleFunc("/api/shadow/diff", shadowPilotHandler)
	http.HandleFunc("/system/status", submoduleStatusHandler)
	http.HandleFunc("/api/shadow/autofix", ciAutoFixHandler)
	http.HandleFunc("/api/queue/telemetry", queueTelemetryHandler)
	http.HandleFunc("/api/tasks/route", taskRouterHandler)
	http.HandleFunc("/api/conflicts/resolve", conflictResolutionHandler)
	http.HandleFunc("/api/system/drift", driftDetectionHandler)
	http.HandleFunc("/api/system/prune", pruneSubmodulesHandler)
	http.HandleFunc("/api/system/dashboard", generateDashboardHandler)

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
