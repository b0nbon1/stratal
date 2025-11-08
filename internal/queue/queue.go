package queue

import "time"

type TaskQueue interface {
	Enqueue(jobRunID string) error
	Dequeue(block time.Duration) (string, map[string]interface{}, error)
	Ack(msgID string) error
	RetryOrDLQ(msgID string, values map[string]interface{}) error
	MoveToDeadLetter(values map[string]interface{}) error
	ReclaimStuckJobs(idleTimeout time.Duration)
}
