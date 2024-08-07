package main

import (
	"context"
	"distributed-task-queue/internal/config"
	"distributed-task-queue/internal/handlers"
	"distributed-task-queue/internal/models"
	"distributed-task-queue/internal/queue"
	"distributed-task-queue/internal/worker"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"distributed-task-queue/internal/scheduler"
)

func main() {
    cfg := config.LoadConfig()

    redisClient := redis.NewClient(&redis.Options{
        Addr: cfg.RedisAddr,
    })

    queue := queue.NewRedisQueue(redisClient, "job_queue")
    worker := worker.NewWorker(queue, cfg.WorkerCount)

    worker.RegisterHandler("email", handlers.HandleEmailJob)
    worker.RegisterHandler("image_resize", handlers.HandleImageResizeJob)
    worker.RegisterHandler("data_processing", handlers.HandleDataProcessingJob)

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.WorkerTimeout)*time.Second)
    defer cancel()

    worker.Start(ctx)

    scheduler := scheduler.NewScheduler(queue)

    // Example: Schedule a job to run every minute
    emailJob := &models.Job{
        Type: "email",
        Payload: json.RawMessage(`{
            "to": "user@example.com",
            "subject": "Scheduled Email",
            "body": "This is a scheduled email."
        }`),
    }
    err := scheduler.AddJob("* * * * *", emailJob)
    if err != nil {
        log.Fatalf("Failed to schedule job: %v", err)
    }

    scheduler.Start()
    defer scheduler.Stop()

    // Set up API server
    http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var job models.Job
        if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        job.ID = uuid.New().String()
        job.Status = "pending"
        job.CreatedAt = time.Now()
        job.UpdatedAt = time.Now()

        if err := queue.Enqueue(r.Context(), &job); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusAccepted)
        json.NewEncoder(w).Encode(map[string]string{"job_id": job.ID})
    })

    log.Printf("Starting API server on %s", cfg.APIAddr)
    go func() {
        if err := http.ListenAndServe(cfg.APIAddr, nil); err != nil {
            log.Fatalf("Failed to start API server: %v", err)
        }
    }()

    select {} // Keep the main goroutine running
}
