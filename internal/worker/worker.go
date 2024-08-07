package worker 


import (
	"context"
	"distributed-task-queue/internal/queue"
	"distributed-task-queue/internal/models"
	"log"
)

type Worker struct {
    queue       queue.Queue
    handlers    map[string]JobHandler
    concurrency int
}

type JobHandler func(context.Context, * models.Job) error

func NewWorker(queue queue.Queue, concurrency int) *Worker {
    return &Worker{
        queue:       queue,
        handlers:    make(map[string]JobHandler),
        concurrency: concurrency,
    }
}

func (w *Worker) RegisterHandler(jobType string, handler JobHandler) {
    w.handlers[jobType] = handler
}

func (w *Worker) Start(ctx context.Context) {
    for i := 0; i < w.concurrency; i++ {
        go w.work(ctx)
    }
}

func (w *Worker) work(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            job, err := w.queue.Dequeue(ctx)
            if err != nil {
                log.Printf("Error dequeuing job: %v", err)
                continue
            }
            if handler, ok := w.handlers[job.Type]; ok {
                job.Status = "processing"
                w.queue.UpdateJob(ctx, job)
                err = handler(ctx, job)
                if err != nil {
                    job.Status = "failed"
                    log.Printf("Error processing job %s: %v", job.ID, err)
                } else {
                    job.Status = "completed"
                }
                w.queue.UpdateJob(ctx, job)
            } else {
                log.Printf("No handler for job type: %s", job.Type)
            }
        }
    }
}
