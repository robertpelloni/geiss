package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Register handlers
	http.HandleFunc("/api/shadow/diff", shadowPilotHandler)
	http.HandleFunc("/system/status", submoduleStatusHandler)
	http.HandleFunc("/api/shadow/autofix", ciAutoFixHandler)

	// Start background daemons
	startShadowDaemon()

	fmt.Println("Jules Autopilot Go Backend starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
