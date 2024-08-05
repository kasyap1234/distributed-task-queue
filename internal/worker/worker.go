package worker

import (
	"context"
	"distributed-task-queue/internal/logger"
	"distributed-task-queue/internal/models"
	"distributed-task-queue/internal/queue"
	"fmt"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
)
type Worker struct {
	queue *queue.RedisQueue
	maxRetries int
	rateLimit int 
	numWorkers int 

}

func NewWorker(q *queue.RedisQueue,maxRetries,rateLimit int) *Worker {
	return &Worker{
		queue: q,
		maxRetries: maxRetries,
		rateLimit: rateLimit,
	}
}
func (w *Worker) Start (ctx context.Context,wg *sync.WaitGroup){
	defer wg.Done()
	for  i :=0; i< w.numWorkers ; i++ {
		wg.Add(1); 
		go w.runWorker(ctx,wg)
	}
}
func (w *Worker) runWorker (ctx context.Context,wg *sync.WaitGroup){
	defer wg.Done()
	rateLimiter :=time.Tick(time.Second /time.Duration(w.rateLimit))
	for {
		select {
		case <-ctx.Done():
			logger.Logger.Info("Worker shutting Down"); 
			return 
		case <-rateLimiter:
			job,err :=w.queue.Dequeue(); 

			// if queue is empty then wait for a second 
			if err ==queue.ErrQueueEmpty{
				time.Sleep(time.Second); 
				continue
			}
			if err!=nil {
				logger.Logger.Error("Error dequeueing job",zap.Error())
				time.Sleep(time.Second); 
				continue 

			}
			if err :=w.processJob(job); err !=nil {
				
			}
		}
	}
}
