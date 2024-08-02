package worker 
import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"distributed-task-queue/internal/models"
	"distributed-task-queue/internal/queue"
	"distributed-task-queue/internal/logger"

)
type Worker struct {
	queue *queue.RedisQueue
	maxRetries int 
}
func NewWorker 