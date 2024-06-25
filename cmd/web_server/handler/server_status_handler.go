package handler

import (
	"fmt"
	"net/http"
	"time"
)

// ServerStatusHandler a Handler for getting the server status.
type ServerStatusHandler struct {
	startTime time.Time
	version   string
}

// NewServerStatusHandler creates a new ServerStatusHandler instance.
func NewServerStatusHandler(version string) *ServerStatusHandler {
	return &ServerStatusHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// Handle is responsible for handling the request and response for the server status.
func (ssh *ServerStatusHandler) Handle(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(ssh.startTime).String()

	fmt.Fprintf(w, "Server Status:\n")
	fmt.Fprintf(w, "Version: %s\n", ssh.version)
	fmt.Fprintf(w, "Uptime: %s\n", uptime)
	fmt.Fprintf(w, "Server Time: %s\n", time.Now().Format(time.RFC3339))
}
