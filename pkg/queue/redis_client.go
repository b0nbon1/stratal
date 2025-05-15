package queue

import (
	"context"
	"fmt"

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

func (rq *RedisQueue) Enqueue(job []byte) {
    err := rq.client.RPush(rq.ctx, rq.key, job).Err()
    if err != nil {
        panic(err)
    }
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

func (rq *RedisQueue) XReadGeneric(lastID string, block int64, mapper func(map[string]interface{}) (any, error)) ([]any, error) {
    args := &redis.XReadArgs{
        Streams: []string{rq.key, lastID},
        Count:   1,
        Block:  0,
    }

    streams, err := rq.client.XRead(rq.ctx, args).Result()
    if err != nil {
        return nil, err
    }

    if len(streams) == 0 || len(streams[0].Messages) == 0 {
        return nil, nil
    }

    tasks := make([]any, 0, len(streams[0].Messages))

    for _, msg := range streams[0].Messages {
        task, err := mapper(msg.Values)

		if err != nil {
            return nil, err
        }
        taskID := msg.ID
        taskMap, ok := task.(map[string]any)
        if !ok {
            return nil, fmt.Errorf("unexpected type %T", task)
        }
        taskMap["id"] = taskID
        taskMap["data"] = task
        tasks = append(tasks, taskMap)
    }
    return tasks, nil
}

func (rq *RedisQueue) XDelete(id string) error {
    _, err := rq.client.XDel(rq.ctx, rq.key, id).Result()
    if err != nil {
        return err
    }
    return nil
}

