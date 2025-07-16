package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
    client *redis.Client
    key    string
    ctx    context.Context
}

func NewRedisQueue(addr, password string, db int, key string) *RedisQueue {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    return &RedisQueue{
        client: rdb,
        key:    key,
        ctx:    context.Background(),
    }
}

func (rq *RedisQueue) Enqueue(job_run_id string) error {
    err := rq.client.RPush(rq.ctx, rq.key, job_run_id).Err()
    if err != nil {
        return fmt.Errorf("unable to queue Job_run_id '%s' with error: %w", job_run_id, err)
    }
    return nil
}

func (rq *RedisQueue) Dequeue() (string, error) {
    result, err := rq.client.BLPop(rq.ctx, 1*time.Second, rq.key).Result()
    if err != nil {
        return "", err
    }
    if len(result) < 2 {
        return "", nil
    }
	
    return result[1], nil
}

func (rq *RedisQueue) MoveToDeadLetter(task []byte) {
	dlqKey := rq.key + ":dlq"
	err := rq.client.RPush(rq.ctx, dlqKey, task).Err()
	if err != nil {
		fmt.Println("Error moving task to DLQ:", err)
	}
}
