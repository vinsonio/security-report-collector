package queue

import (
	"sync"
)

// InMemoryQueue is an in-memory queue implementation (not persistent).
type InMemoryQueue struct {
	mutex   sync.Mutex
	items   []*ReportEnvelope
	hashSet map[string]bool
}

// NewInMemoryQueue creates a new in-memory queue.
func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		items:   make([]*ReportEnvelope, 0),
		hashSet: make(map[string]bool),
	}
}

// Enqueue adds a report envelope to the queue.
func (q *InMemoryQueue) Enqueue(envelope *ReportEnvelope) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, envelope)
	q.hashSet[envelope.Hash] = true
	return nil
}

// DequeueN retrieves and removes up to n envelopes from the queue.
func (q *InMemoryQueue) DequeueN(n int) ([]*ReportEnvelope, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return []*ReportEnvelope{}, nil
	}

	count := n
	if count > len(q.items) {
		count = len(q.items)
	}

	result := make([]*ReportEnvelope, count)
	copy(result, q.items[:count])

	// Remove dequeued items and their hashes
	for _, envelope := range result {
		delete(q.hashSet, envelope.Hash)
	}

	q.items = q.items[count:]
	return result, nil
}

// Size returns the number of items in the queue.
func (q *InMemoryQueue) Size() (int, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.items), nil
}

// Contains checks if a hash exists in the queue (for deduplication).
func (q *InMemoryQueue) Contains(hash string) (bool, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.hashSet[hash], nil
}

// Close closes the queue.
func (q *InMemoryQueue) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items = nil
	q.hashSet = nil
	return nil
}
