package scheduler

import (
	"log"

	"github.com/vinsonio/security-report-collector/internal/queue"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// Database is the interface for database operations used by the flusher.
type Database interface {
	Save(reportType string, report types.Report, userAgent string, hash string) error
}

// BatchFlusher is responsible for flushing queued reports to the database.
type BatchFlusher struct {
	queue     queue.Queue
	database  Database
	batchSize int
}

// NewBatchFlusher creates a new batch flusher.
func NewBatchFlusher(q queue.Queue, db Database, batchSize int) *BatchFlusher {
	return &BatchFlusher{
		queue:     q,
		database:  db,
		batchSize: batchSize,
	}
}

// Flush dequeues and persists up to batchSize reports from the queue.
func (f *BatchFlusher) Flush() error {
	envelopes, err := f.queue.DequeueN(f.batchSize)
	if err != nil {
		return err
	}

	if len(envelopes) == 0 {
		// Nothing to flush
		return nil
	}

	log.Printf("Flushing %d reports to database", len(envelopes))

	successCount := 0
	for _, envelope := range envelopes {
		err := f.database.Save(envelope.Type, envelope.Report, envelope.UserAgent, envelope.Hash)
		if err != nil {
			log.Printf("Failed to save report (hash: %s, type: %s): %v", envelope.Hash, envelope.Type, err)
			// Continue with other reports - don't fail the entire batch
		} else {
			successCount++
		}
	}

	log.Printf("Successfully flushed %d/%d reports", successCount, len(envelopes))
	return nil
}
