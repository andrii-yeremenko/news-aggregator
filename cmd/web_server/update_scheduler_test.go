package web_server

import (
	"bytes"
	"log"
	"testing"
	"time"
)

// MockManager implements the Manager interface for testing purposes.
type MockManager struct {
	UpdateAllSourcesFunc func() error
}

// UpdateAllSources calls the mocked function.
func (m *MockManager) UpdateAllSources() error {
	return m.UpdateAllSourcesFunc()
}

// TestUpdateScheduler_Start tests the Start method of the UpdateScheduler.
func TestUpdateScheduler_LogMessages(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	mockManager := &MockManager{
		UpdateAllSourcesFunc: func() error {
			return nil
		},
	}

	timeout := time.Millisecond * 100
	scheduler := NewUpdateScheduler(mockManager, timeout)

	done := make(chan struct{})
	go func() {
		defer close(done)
		scheduler.Start()
	}()

	time.Sleep(timeout * 3)

	expectedMessages := []string{
		"Updating resources...",
		"Resources updated at",
	}

	for _, expectedMsg := range expectedMessages {
		if !containsLogMessage(&buf, expectedMsg) {
			t.Errorf("Expected log message '%s' was not found", expectedMsg)
		}
	}
}

// containsLogMessage checks if the expected message is contained in the log buffer
func containsLogMessage(buf *bytes.Buffer, expectedMsg string) bool {
	return bytes.Contains(buf.Bytes(), []byte(expectedMsg))
}
