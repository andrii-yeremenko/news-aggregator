package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestStatusHandler_New tests the New function of StatusHandler.
func TestStatusHandler_New(t *testing.T) {

	version := "1.0.0"
	handler := NewStatusHandler(version)

	if handler == nil {
		t.Fatalf("expected non-nil handler")
	}

	if handler.version != version {
		t.Errorf("expected version %s, got %s", version, handler.version)
	}

	if handler.startTime.IsZero() {
		t.Errorf("expected non-zero start time")
	}
}

// TestStatusHandler_Handle tests the Handle function of StatusHandler.
func TestStatusHandler_Handle(t *testing.T) {
	startTime := time.Now().Add(-1 * time.Hour)
	version := "1.0.0"
	handler := &StatusHandler{
		startTime: startTime,
		version:   version,
	}

	req := httptest.NewRequest("GET", "http://example.com/status", nil)
	rr := httptest.NewRecorder()

	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedVersion := "Version: 1.0.0\n"
	if !strings.Contains(rr.Body.String(), expectedVersion) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedVersion)
	}

	if !strings.Contains(rr.Body.String(), "Uptime: 1h") {
		t.Errorf("handler returned unexpected body: got %v want Uptime to contain '1h'",
			rr.Body.String())
	}

	expectedTime := "Server Time: " // We can't predict the exact time, so just check the prefix.
	if !strings.Contains(rr.Body.String(), expectedTime) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedTime)
	}
}
