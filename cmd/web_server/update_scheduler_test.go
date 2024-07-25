package web_server

import (
	"bytes"
	"errors"
	"log"
	"news-aggregator/cmd/web_server/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

// TestUpdateScheduler_LogMessages tests the Start method of the UpdateScheduler for positive cases.
func TestUpdateScheduler_LogMessages(t *testing.T) {
	var buf bytes.Buffer
	originalLogOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalLogOutput)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mocks.NewMockManager(ctrl)
	mockManager.EXPECT().UpdateAllSources().Return(nil).AnyTimes()

	timeout := time.Millisecond * 100
	scheduler := NewUpdateScheduler(mockManager, timeout)

	done := make(chan struct{})
	go func() {
		defer close(done)
		scheduler.Start()
	}()

	time.Sleep(timeout * 3)
	scheduler.Stop()

	<-done

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

// TestUpdateScheduler_ErrorLog tests the Start method of the UpdateScheduler for negative cases.
func TestUpdateScheduler_ErrorLog(t *testing.T) {
	var buf bytes.Buffer
	originalLogOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalLogOutput)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mocks.NewMockManager(ctrl)
	mockManager.EXPECT().UpdateAllSources().Return(errors.New("update failed")).AnyTimes()

	timeout := time.Millisecond * 100
	scheduler := NewUpdateScheduler(mockManager, timeout)

	done := make(chan struct{})
	go func() {
		defer close(done)
		scheduler.Start()
	}()

	time.Sleep(timeout * 3)
	scheduler.Stop()

	<-done

	expectedErrorMessage := "Failed to update resources: update failed"

	if !containsLogMessage(&buf, expectedErrorMessage) {
		t.Errorf("Expected log message '%s' was not found", expectedErrorMessage)
	}
}

// containsLogMessage checks if the expected message is contained in the log buffer.
func containsLogMessage(buf *bytes.Buffer, expectedMsg string) bool {
	return bytes.Contains(buf.Bytes(), []byte(expectedMsg))
}
