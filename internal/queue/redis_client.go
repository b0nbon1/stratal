// package queue

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/b0nbon1/stratal/internal/config"
// 	"github.com/redis/go-redis/v9"
// )

// type RedisQueue struct {
// 	client *redis.Client
// 	key    string
// 	ctx    context.Context
// }

// func NewRedisQueue(cfg *config.Config, key string) *RedisQueue {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     cfg.Redis.Addr,
// 		Password: cfg.Redis.Password,
// 		DB:       cfg.Redis.DB,
// 	})
// 	return &RedisQueue{
// 		client: rdb,
// 		key:    key,
// 		ctx:    context.Background(),
// 	}
// }

// func (rq *RedisQueue) Enqueue(job_run_id string) error {
// 	err := rq.client.RPush(rq.ctx, rq.key, job_run_id).Err()
// 	if err != nil {
// 		return fmt.Errorf("unable to queue Job_run_id '%s' with error: %w", job_run_id, err)
// 	}
// 	return nil
// }

// func (rq *RedisQueue) Dequeue() (string, error) {
// 	result, err := rq.client.BLPop(rq.ctx, 1*time.Second, rq.key).Result()
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(result) < 2 {
// 		return "", nil
// 	}

// 	return result[1], nil
// }

// func (rq *RedisQueue) MoveToDeadLetter(task []byte) {
// 	dlqKey := rq.key + ":dlq"
// 	err := rq.client.RPush(rq.ctx, dlqKey, task).Err()
// 	if err != nil {
// 		fmt.Println("Error moving task to DLQ:", err)
// 	}
// }

package queue

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/b0nbon1/stratal/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client    *redis.Client
	stream    string
	group     string
	consumer  string
	ctx       context.Context
	maxRetries int
}

func NewRedisQueue(cfg *config.Config, stream, group string, maxRetries int) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()

	// Try to create consumer group (ignore "already exists" errors)
	_ = rdb.XGroupCreateMkStream(ctx, stream, group, "$").Err()

	// dynamic consumer name: use hostname + PID
	hostname, _ := os.Hostname()
	consumer := fmt.Sprintf("%s-%d", hostname, os.Getpid())

	return &RedisQueue{
		client:    rdb,
		stream:    stream,
		group:     group,
		consumer:  consumer,
		ctx:       ctx,
		maxRetries: maxRetries,
	}
}

// Enqueue pushes a job with retry count = 0
func (rq *RedisQueue) Enqueue(jobRunID string) error {
	err := rq.client.XAdd(rq.ctx, &redis.XAddArgs{
		Stream: rq.stream,
		Values: map[string]interface{}{
			"job_run_id":  jobRunID,
			"retry_count": 0,
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("unable to enqueue Job_run_id '%s': %w", jobRunID, err)
	}
	return nil
}

// Dequeue gets a job (blocking up to `block`) and returns ID + values
func (rq *RedisQueue) Dequeue(block time.Duration) (string, map[string]interface{}, error) {
	streams, err := rq.client.XReadGroup(rq.ctx, &redis.XReadGroupArgs{
		Group:    rq.group,
		Consumer: rq.consumer,
		Streams:  []string{rq.stream, ">"},
		Count:    1,
		Block:    block,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			return "", nil, nil
		}
		return "", nil, err
	}

	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return "", nil, nil
	}

	msg := streams[0].Messages[0]
	return msg.ID, msg.Values, nil
}

// Ack acknowledges successful processing
func (rq *RedisQueue) Ack(msgID string) error {
	return rq.client.XAck(rq.ctx, rq.stream, rq.group, msgID).Err()
}

// RetryOrDLQ handles failed jobs
func (rq *RedisQueue) RetryOrDLQ(msgID string, values map[string]interface{}) error {
	// read current retry count
	retries := 0
	if val, ok := values["retry_count"]; ok {
		switch v := val.(type) {
		case int64:
			retries = int(v)
		case int:
			retries = v
		case string:
			// Redis might return string
			var tmp int
			fmt.Sscanf(v, "%d", &tmp)
			retries = tmp
		}
	}

	if retries < rq.maxRetries {
		// re-enqueue with incremented retry count
		values["retry_count"] = retries + 1
		return rq.client.XAdd(rq.ctx, &redis.XAddArgs{
			Stream: rq.stream,
			Values: values,
		}).Err()
	}

	// move to dead-letter
	return rq.MoveToDeadLetter(values)
}

// MoveToDeadLetter sends job to DLQ
func (rq *RedisQueue) MoveToDeadLetter(values map[string]interface{}) error {
	dlqKey := rq.stream + ":dlq"
	return rq.client.XAdd(rq.ctx, &redis.XAddArgs{
		Stream: dlqKey,
		Values: values,
	}).Err()
}

// ReclaimStuckJobs finds unacked jobs idle > idleTimeout and reclaims them
func (rq *RedisQueue) ReclaimStuckJobs(idleTimeout time.Duration) {
	pending, err := rq.client.XPendingExt(rq.ctx, &redis.XPendingExtArgs{
		Stream:   rq.stream,
		Group:    rq.group,
		Idle:     idleTimeout,
		Count:    10,
		Start:    "-",
		End:      "+",
	}).Result()

	if err != nil {
		return
	}

	for _, p := range pending {
		_, err := rq.client.XClaim(rq.ctx, &redis.XClaimArgs{
			Stream:   rq.stream,
			Group:    rq.group,
			Consumer: rq.consumer,
			MinIdle:  idleTimeout,
			Messages: []string{p.ID},
		}).Result()
		if err != nil {
			fmt.Println("XCLAIM error:", err)
		}
	}
}

