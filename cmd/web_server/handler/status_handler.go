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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := time.Since(ssh.startTime).String()
	currentTime := time.Now().Format(time.RFC3339)

	if _, err := fmt.Fprintf(w, "Server Status:\n"); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Version: %s\n", ssh.version); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Uptime: %s\n", uptime); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Server Time: %s\n", currentTime); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
