package queue 

import (
    "context"
    "encoding/json"
    "fmt"
    
	"distributed-task-queue/internal/models"
    "github.com/go-redis/redis/v8"
)

type Queue interface {
    Enqueue(ctx context.Context, job *models.Job) error
    Dequeue(ctx context.Context) (*models.Job, error)
    UpdateJob(ctx context.Context, job *models.Job) error
}

type RedisQueue struct {
    client *redis.Client
    key    string
}

func NewRedisQueue(client *redis.Client, key string) *RedisQueue {
    return &RedisQueue{client: client, key: key}
}

func (q *RedisQueue) Enqueue(ctx context.Context, job *models.Job) error {
    jobJSON, err := json.Marshal(job)
    if err != nil {
        return fmt.Errorf("failed to marshal job: %w", err)
    }
    return q.client.RPush(ctx, q.key, jobJSON).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (*models.Job, error) {
    result, err := q.client.BLPop(ctx, 0, q.key).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to pop job from queue: %w", err)
    }
    var job models.Job
    err = json.Unmarshal([]byte(result[1]), &job)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal job: %w", err)
    }
    return &job, nil
}

func (q *RedisQueue) UpdateJob(ctx context.Context, job *models.Job) error {
    jobJSON, err := json.Marshal(job)
    if err != nil {
        return fmt.Errorf("failed to marshal job: %w", err)
    }
    return q.client.Set(ctx, fmt.Sprintf("job:%s", job.ID), jobJSON, 0).Err()
}
