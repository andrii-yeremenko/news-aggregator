package handler

import (
	"fmt"
	"net/http"
	"time"
)

// StatusHandler a Handler for getting the server status.
type StatusHandler struct {
	startTime time.Time
	version   string
}

// NewStatusHandler creates a new StatusHandler instance.
func NewStatusHandler(version string) *StatusHandler {
	return &StatusHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// Handle is responsible for handling the request and response for the server status.
func (ssh *StatusHandler) Handle(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(ssh.startTime).String()

	fmt.Fprintf(w, "Server Status:\n")
	fmt.Fprintf(w, "Version: %s\n", ssh.version)
	fmt.Fprintf(w, "Uptime: %s\n", uptime)
	fmt.Fprintf(w, "Server Time: %s\n", time.Now().Format(time.RFC3339))
}
