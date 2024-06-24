package handler

import (
	"fmt"
	"net/http"
	"time"
)

type ServerStatusHandler struct {
	startTime time.Time
	version   string
}

func NewServerStatusHandler(version string) *ServerStatusHandler {
	return &ServerStatusHandler{
		startTime: time.Now(),
		version:   version,
	}
}

func (ssh *ServerStatusHandler) Handle(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(ssh.startTime).String()

	fmt.Fprintf(w, "Server Status:\n")
	fmt.Fprintf(w, "Version: %s\n", ssh.version)
	fmt.Fprintf(w, "Uptime: %s\n", uptime)
	fmt.Fprintf(w, "Server Time: %s\n", time.Now().Format(time.RFC3339))
}
