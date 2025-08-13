package scheduler

import (
	"log"
	"time"
)

// Scheduler runs a function at a fixed interval until the stop channel is closed.
func Scheduler(interval time.Duration, stop <-chan struct{}, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				log.Printf("scheduled job error: %v", err)
			}
		case <-stop:
			log.Println("scheduler stopped")
			return
		}
	}
}
