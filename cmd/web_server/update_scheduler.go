package web_server

import (
	"log"
	"time"
)

// UpdateScheduler schedules updates based on a timeout duration.
type UpdateScheduler struct {
	manager Manager
	timeout time.Duration
	stop    chan struct{} // Channel to signal stopping the scheduler
}

// NewUpdateScheduler creates a new UpdateScheduler instance.
func NewUpdateScheduler(m Manager, timeout time.Duration) *UpdateScheduler {
	return &UpdateScheduler{
		manager: m,
		timeout: timeout,
		stop:    make(chan struct{}),
	}
}

// Start starts the update scheduling process in a separate goroutine.
func (s *UpdateScheduler) Start() {
	log.Printf("Starting update scheduler with timeout %s ...\n", s.timeout.String())

	go func() {
		ticker := time.NewTicker(s.timeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Updating resources...")
				err := s.manager.UpdateAllSources()
				if err != nil {
					log.Printf("Failed to update resources: %v", err)
				}
				t := time.Now().Format("2006-01-02 15:04:05")
				log.Println("Resources updated at", t)

			case <-s.stop:
				log.Println("Stopping update scheduler...")
				return
			}
		}
	}()
}

// Stop stops the update scheduler.
func (s *UpdateScheduler) Stop() {
	close(s.stop)
}
