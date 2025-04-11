package queue

import (
	"context"
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

func (rq *RedisQueue) Enqueue(task string) error {
    return rq.client.RPush(rq.ctx, rq.key, task).Err()
}

func (rq *RedisQueue) Dequeue() (string, error) {
    result, err := rq.client.BLPop(rq.ctx, 0, rq.key).Result()
    if err != nil {
        return "", err
    }
    if len(result) < 2 {
        return "", nil
    }
    return result[1], nil
}

func (rq *RedisQueue) XRead(lastID string, block int64) ([]string, error) {
    args := &redis.XReadArgs{
        Streams: []string{rq.key, lastID},
        Count:   1,
        Block:   time.Duration(block) * time.Millisecond,
    }

    streams, err := rq.client.XRead(rq.ctx, args).Result()
    if err != nil {
        return nil, err
    }

    if len(streams) == 0 || len(streams[0].Messages) == 0 {
        return nil, nil
    }

    tasks := make([]string, 0, len(streams[0].Messages))

    for _, msg := range streams[0].Messages {
        if taskVal, ok := msg.Values["task"]; ok {
            if taskStr, valid := taskVal.(string); valid {
                tasks = append(tasks, taskStr)
            }
        }
    }
    return tasks, nil
}

