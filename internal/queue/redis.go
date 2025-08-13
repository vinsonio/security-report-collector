package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisQueue is a Redis-based queue implementation.
type RedisQueue struct {
	client   *redis.Client
	queueKey string
	hashKey  string
	ctx      context.Context
}

// NewRedisQueue creates a new Redis queue.
func NewRedisQueue(addr, password string, db int, queueName string) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return &RedisQueue{
		client:   client,
		queueKey: "queue:" + queueName,
		hashKey:  "queue:" + queueName + ":hashes",
		ctx:      ctx,
	}, nil

}

// Enqueue adds a report envelope to the queue.
func (q *RedisQueue) Enqueue(envelope *ReportEnvelope) error {
	data, err := MarshalEnvelope(envelope)
	if err != nil {
		return err
	}

	// Use a pipeline for atomicity
	pipe := q.client.TxPipeline()

	// Add to list (queue)
	pipe.LPush(q.ctx, q.queueKey, data)

	// Track hash for deduplication
	pipe.SAdd(q.ctx, q.hashKey, envelope.Hash)

	_, err = pipe.Exec(q.ctx)

	return err

}

// DequeueN retrieves and removes up to n envelopes from the queue.
func (q *RedisQueue) DequeueN(n int) ([]*ReportEnvelope, error) {
	var envelopes []*ReportEnvelope

	for i := 0; i < n; i++ {
		// Pop from the right (FIFO)
		result, err := q.client.RPop(q.ctx, q.queueKey).Result()
		if err == redis.Nil {
			// Queue is empty
			break
		}
		if err != nil {
			return nil, err
		}

		envelope, err := UnmarshalEnvelope([]byte(result))
		if err != nil {
			return nil, err
		}

		envelopes = append(envelopes, envelope)

		// Remove hash from set
		q.client.SRem(q.ctx, q.hashKey, envelope.Hash)
	}

	return envelopes, nil
}

// Size returns the approximate number of items in the queue.
func (q *RedisQueue) Size() (int, error) {
	size, err := q.client.LLen(q.ctx, q.queueKey).Result()
	return int(size), err
}

// Contains checks if a hash exists in the queue (for deduplication).
func (q *RedisQueue) Contains(hash string) (bool, error) {
	return q.client.SIsMember(q.ctx, q.hashKey, hash).Result()
}

// Close closes the queue.
func (q *RedisQueue) Close() error {
	return q.client.Close()
}
