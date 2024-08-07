package scheduler

import (
	"context"
	"distributed-task-queue/internal/models"
	"distributed-task-queue/internal/queue"
	

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
    cron  *cron.Cron
    queue queue.RedisQueue
}

func NewScheduler(queue *queue.RedisQueue) *Scheduler {
    return &Scheduler{
        cron:  cron.New(),
        queue: *queue,
    }
}

func (s *Scheduler) AddJob(spec string, job *models.Job) error {
    _, err := s.cron.AddFunc(spec, func() {
        ctx := context.Background()
        s.queue.Enqueue(ctx, job)
    })
    return err
}

func (s *Scheduler) Start() {
    s.cron.Start()
}

func (s *Scheduler) Stop() {
    s.cron.Stop()
}
